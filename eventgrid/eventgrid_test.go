package eventgrid

// MyEventSpecificData demonstrates how one can declare an aribtrary type to be
// sent and received as part of an EventGrid call.
type MyEventSpecificData struct {
	Field1 string `json:"field1"`
	Field2 string `json:"field2"`
	Field3 string `json:"field3"`
}

func ExampleCreateTopic() {
	// Output:
}
