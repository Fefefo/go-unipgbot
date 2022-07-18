package cache

import (
	"fmt"
	"main/enums/states"
	"math/rand"
	"os"
	"testing"
)

func TestState(t *testing.T) {
	var err error

	rn := rand.Int63()
	// Test caso che non esiste (deve restituire l'errore stateless user)
	if _, err = State(rn); err.Error() != "stateless user" {
		t.Fail()
	}

	// Inserisco uno stato per l'id 1234 (deve funzionare)
	if _, err = State(rn, states.Home); err != nil {
		t.Fail()
	}

	// Testo di nuovo il caso ora che l'id esiste
	if _, err = State(rn); err != nil {
		t.Fail()
	}

	// Cancello la cartella che si è creata
	if err = os.RemoveAll("cache"); err != nil {
		fmt.Println("Test passato con successo, ma la cartella non è stata cancellata")
	}
}
