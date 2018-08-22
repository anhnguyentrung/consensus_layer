package serializer

import (
	"testing"
	"fmt"
)

type User struct {
	Id 		int32
	Name 	string
}

func newUser(id int32, name string) User {
	return User {
		Id: 	id,
		Name: 	name,
	}
}

func TestStructSerializer(t *testing.T) {
	user1 := newUser(1, "user1")
	//fmt.Println(user1)
	buf, _ := MarshalBinary(user1)
	user2 := newUser(0, "")
	UnmarshalBinary(buf, &user2)
	fmt.Println(user2)
}

func TestMapSerializer(t *testing.T) {
	users := make(map[User]bool, 0)
	user1 := newUser(1, "user1")
	users[user1] = true
	buf, _ := MarshalBinary(users)
	users2 := make(map[User]bool, 0)
	UnmarshalBinary(buf, &users2)
	fmt.Println(users2)
}
