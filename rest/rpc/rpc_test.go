package rpc

import (
	"encoding/json"
	"testing"
)

func TestJob_UnmarshalJSON(t *testing.T) {
	input := []byte(`{
	"headers": {
		"connector": 42
	},
	"body": "dGhpcyBpcyB0aGUgc29uZyB0aGF0IG5ldmVyIGVuZHM="
}`)
	job := &Job{}
	err := json.Unmarshal(input, job)
	if err != nil {
		t.Fatal(err)
	}
	if string(job.Body) != "this is the song that never ends" {
		t.Log(job.Body)
		t.Fatalf("got '%s' want 'this is the song that never ends'", string(job.Body))
	}
}

func TestJob_MarshalJSON(t *testing.T) {
	serverSent := &Job{
		Headers: Headers{Connector: 42},
		Body:    Body("this is the song that never ends"),
	}
	input, err := json.Marshal(serverSent)
	if err != nil {
		t.Fatal(err)
	}
	job := &Job{}
	err = json.Unmarshal(input, job)
	if err != nil {
		t.Fatal(err)
	}
	if string(job.Body) != "this is the song that never ends" {
		t.Log(job.Body)
		t.Fatalf("got '%s' want 'this is the song that never ends'", string(job.Body))
	}
}

func TestUnmarshal(t *testing.T) {
	target := make([]byte, 0)
	incoming := []byte("lol sup my guy")
	Unmarshal(incoming, &target)
	t.Logf("come on '%s'", string(target))
}

func Unmarshal(data []byte, v interface{}) error {
	b, _ := v.(*[]byte)
	*b = data
	return nil
}

func TestOutOfBound(t *testing.T) {
	b := &Body{}
	json.Unmarshal([]byte(`"CgoxPz81WFdSRVBF"`), b)
	t.Log(string(*b))
}

func TestEmtpy(t *testing.T) {
	b := &Body{}
	json.Unmarshal([]byte(`""`), b)
	t.Log(string(*b))
}
