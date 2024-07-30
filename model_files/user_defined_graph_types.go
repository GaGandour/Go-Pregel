package graph_package

/*
You can change the VertexIdType to any type you want.
But please, don't get to crazy with it.
Strings always work well.
*/
type VertexIdType string
type EdgeIdType string

type VertexValue struct {
	/*
		This is the value that the vertex holds.
		You can put how many fields you want here, but remember:

		This value is going to be read/written from/to the disk.
		You should use serializable types AND the name of the field
		must begin with a capital letter.
	*/
}

type EdgeValue struct {
	/*
		This is the value that the edge holds.
		You can put how many fields you want here, but remember:

		This value is going to be read/written from/to the disk.
		You should use serializable types AND the name of the field
		must begin with a capital letter.
	*/
}

type PregelMessage struct {
	/*
		This is the pregel message struct.

		This will be passed through RPC calls, so you should use
		serializable types AND the name of the field must begin
		with a capital letter so that you can access the fields
		in the Compute method.
	*/
}
