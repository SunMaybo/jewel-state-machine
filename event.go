package jewel_state_machine

type Event struct {
	Src  string
	Event string


	// Err is an optional error that can be returned from a callback.
	Err error

	// Args is a optinal list of arguments passed to the callback.
	Args []interface{}

}
