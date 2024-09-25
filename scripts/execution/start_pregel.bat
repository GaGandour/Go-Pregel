@echo off

REM Get arguments

set DEBUG=false
set FAILURE_STEP=-1
set CHECKPOINT_FREQUENCY=-1
set TEST=false

:parseArgs
if "%~1"=="" goto endArgs
  if "%~1"=="-h" goto help
  if "%~1"=="--help" goto help

  for /f "tokens=1,2 delims==" %%a in ("%~1") do (
    if "%%a"=="-num_workers" set NUM_WORKERS=%%b
    if "%%a"=="-graph_file" set GRAPH_FILE=%%b
    if "%%a"=="-failure_step" set FAILURE_STEP=%%b
    if "%%a"=="-checkpoint_frequency" set CHECKPOINT_FREQUENCY=%%b
    if "%~1"=="-debug" set DEBUG=true
    if "%~1"=="-test" set TEST=true
  )
shift
goto parseArgs

:help
echo Usage: start_docker.bat -num_workers=<number of workers> -graph_file=<graph input file>
echo The graph file is the relative path from the graphs folder.
echo Optional arguments:
echo   -debug: Run in debug mode. This makes the pregel program to register the graph state in every superstep.
echo   -failure_step=<step number>: Simulate a failure in one of the workers at the specified step number. The worker will not be restarted and the computation will continue.
echo   -test: Run the program in test mode. This will not open the graph visualization.
echo.
echo Example 1: start_docker.bat -num_workers=3 -graph_file=graph1.json
echo Example 2: start_docker.bat -num_workers=3 -graph_file=graph1.json -failure_step=5
echo Example 3: start_docker.bat -num_workers=3 -graph_file=graph1.json -debug
exit /b

:endArgs

REM Check required arguments
if "%NUM_WORKERS%"=="" (
  echo Missing Number of Workers. Run start_docker.bat with -h or --help for more information on the necessary arguments.
  exit /b 1
)
if "%GRAPH_FILE%"=="" (
  echo Missing Graph File. Run start_docker.bat with -h or --help for more information on the necessary arguments.
  exit /b 1
)

call build_docker_image.bat

echo Cleaning outputs from other pregel runs...
call clean_outputs.bat
echo Finished cleaning outputs from other pregel runs

cd ..

if "%DEBUG%"=="true" (
  echo Running in debug mode
  python python-scripts\write_docker_compose.py --num_workers=%NUM_WORKERS% --graph_file=%GRAPH_FILE% --failure_step=%FAILURE_STEP% --checkpoint_frequency=%CHECKPOINT_FREQUENCY% --debug > ..\docker-compose.yml
) else (
  python python-scripts\write_docker_compose.py --num_workers=%NUM_WORKERS% --graph_file=%GRAPH_FILE% --failure_step=%FAILURE_STEP% --checkpoint_frequency=%CHECKPOINT_FREQUENCY% > ..\docker-compose.yml
)

cd ..

REM Create the folder structure for the output graph
for /f "tokens=*" %%i in ("%GRAPH_FILE%") do (
  set graph_output_folder_structure=%%~dpi
)
mkdir "src\output_graphs\%graph_output_folder_structure%"
echo Created folder structure: src\output_graphs\%graph_output_folder_structure%

docker-compose -f docker-compose.yml up -d
echo Starting Pregel with %NUM_WORKERS% workers on file %GRAPH_FILE%
docker attach pregel-master
echo Stopping Pregel containers
cd scripts\execution
call stop_docker_containers.bat

if "%TEST%"=="false" (
  cd ..\..
  call venv\Scripts\activate
  cd visualization
  python draw_graph.py --output_file=%GRAPH_FILE%
  deactivate
  start graph.html
  cd ..
)
