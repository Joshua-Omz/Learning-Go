package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
)

// ---------------------------------------------------------
// 1. Data Structure
// ---------------------------------------------------------
type Note struct {
	ID      string `json:"id"`
	Content string `json:"content"`
}

type QuoteResponse struct {
    Quote  string `json:"quote"`
    Author string `json:"author"`
}

type EnhancedNote struct{
	ID      string `json:"id"`
	Content string `json:"content"`
	Quote QuoteResponse `json:"qoute"`
}

// In-memory database
var db = make(map[string]Note)

const dbFile = "data.json"

// ---------------------------------------------------------
// 2. The Middleware (Security)
// ---------------------------------------------------------
func requireAPIKey(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// TODO: Check the custom HTTP Header "X-API-Key"
		apiKey := r.Header.Get("X-API-Key")

		if apiKey == "" {
			http.Error(w, "missing or invalid api key ", http.StatusUnauthorized)
			return
		}

		next(w, r)
	}
}

// ---------------------------------------------------------
// 3. Endpoint 1: Create Note
// ---------------------------------------------------------
// POST /notes/create
func createNote(w http.ResponseWriter, r *http.Request) {
	var newNote Note
	// TODO: Ensure method is POST
	if r.Method != http.MethodPost {
		http.Error(w, "only accepts post methods", http.StatusBadRequest)
	}
	// TODO: Decode JSON body into Note struct
	err := json.NewDecoder(r.Body).Decode(&newNote)
	if err != nil {
		http.Error(w, "Invalid json ", http.StatusBadRequest)
		return

	}
	// TODO: Add note to map `db`
	db[newNote.ID] = newNote
	// TODO: Save map to `data.json`
	err = saveDBToFile()
	if err != nil {
		http.Error(w, "failed to save state", http.StatusInternalServerError)
		return
	}

	// TODO: Return success message
	fmt.Fprint(w, "successfully created Note")
	fmt.Print("successfully created Note")
}

// ---------------------------------------------------------
// 4. Endpoint 2: Read Note
// ---------------------------------------------------------
// GET /notes/read?id=X
func readNote(w http.ResponseWriter, r *http.Request) {
	// TODO: Ensure method is GET
	if r.Method != http.MethodGet {
		http.Error(w, "you are only trying to retrieve json ", http.StatusBadRequest)
		return
	}
	// TODO: Read `data.json` and load back into `db`
	fileBytes, err := os.ReadFile(dbFile)
	if err != nil {
		http.Error(w, "failed to read json data stored", http.StatusInternalServerError)
		return
	}
	//unzipping the json file from the stored db
	err = json.Unmarshal(fileBytes, &db)
	if err != nil {
		http.Error(w, " failed to parse  data", http.StatusInternalServerError)
		return
	}

	// TODO: Extract `id` from query
	id := r.URL.Query().Get("id")

	// TODO: Return Note as JSON, or 404
	//1. checking if the note exists
	note, exists := db[id]
	if !exists {
		http.Error(w, "the note doesnt exists", http.StatusNotFound)
		return
	}

	// 2.setting the header so the receiver knows how to parse correctly
	w.Header().Set("Content-Type", "application/json")

	//3. Encode the Note struct directly into the reponse writer
	err = json.NewEncoder(w).Encode(note)
	if err != nil {
		http.Error(w, "Error returning Json", http.StatusInternalServerError)
		return
	}
}

// ---------------------------------------------------------
// 5. Endpoint 3: Delete Note
// ---------------------------------------------------------
// DELETE /notes/delete?id=X
func deleteNote(w http.ResponseWriter, r *http.Request) {
	// TODO: Ensure method is DELETE
	if r.Method != http.MethodDelete {
		http.Error(w, "the method is only capable of the delete operation not anything else", http.StatusBadRequest)
	}
	// TODO: Extract `id` from query
	id := r.URL.Query().Get("id")

	// TODO: Delete key from map `db`
	delete(db, id)
	// TODO: Save updated map to `data.json`
	err := saveDBToFile()
	if err != nil {
		http.Error(w,"failed to save state", http.StatusInternalServerError)
		return
	}

	// TODO: Return success message
}

// ---------------------------------------------------------
// 6. Endpoint 4: Enhance Note
// ---------------------------------------------------------
// GET /notes/enhance?id=X
func enhanceNote(w http.ResponseWriter, r *http.Request) {
	// TODO: Retrieve note from `data.json` via ID
	id := r.URL.Query().Get("id")

	fileBytes, err := os.ReadFile(dbFile)

	if err != nil{
		http.Error(w, "failed to read data",http.StatusInternalServerError)
		return
	}
	json.Unmarshal(fileBytes,&db)

	note, exists := db[id]
	if !exists {
		http.Error(w, "the note does not exists",http.StatusInternalServerError)
		return
	}
	
	
	// TODO: Make outbound GET to https://dummyjson.com/quotes/random
	resp, err := http.Get("https://dummyjson.com/quotes/random")
	 if err != nil {
        http.Error(w, "Failed to fetch external quote", http.StatusInternalServerError)
        return
    }
	
	defer resp.Body.Close()

	var randomQuote QuoteResponse
	err = json.NewDecoder(resp.Body).Decode(&randomQuote)
	    if err != nil {
        http.Error(w, "Failed to decode external quote", http.StatusInternalServerError)
        return
    }


	// TODO: Combine original note and random quote
	finalNote := EnhancedNote{
		ID: note.ID,
		Content: note.Content,
		Quote: randomQuote,
	}
	
	 

	// TODO: Return as JSON
		 w.Header().Set("Content-Type", "application/json")
	 err = json.NewEncoder(w).Encode(finalNote)
	   if err != nil {
        http.Error(w, "Failed to encode final specific json", http.StatusInternalServerError)
        return
    }
	err = saveDBToFile()
	if err != nil {
		http.Error(w, "failed to save state", http.StatusInternalServerError)
		return
	}
}

// Helper function
func saveDBToFile() error {
	jsonData, err := json.MarshalIndent(db, "", " ")

	if err != nil {
		fmt.Print("faiied to parse object to json byte slice")
		return err
	}
	err = os.WriteFile(dbFile, jsonData, 0644)
	if err != nil {
		return err
	}
	return nil
}

func main() {
	// Let's create `data.json` if it doesn't exist yet, to prevent crashes early.
	if _, err := os.Stat(dbFile); os.IsNotExist(err) {
		emptyDB := make(map[string]Note)
		b, _ := json.Marshal(emptyDB)
		os.WriteFile(dbFile, b, 0644)
	}

	// ---------------------------------------------------------
	// Routing setup
	// ---------------------------------------------------------
	http.HandleFunc("/notes/create", requireAPIKey(createNote))
	http.HandleFunc("/notes/read", requireAPIKey(readNote))
	http.HandleFunc("/notes/delete", requireAPIKey(deleteNote))
	http.HandleFunc("/notes/enhance", requireAPIKey(enhanceNote))

	fmt.Println("Exam Server starts on port 8080...")
	http.ListenAndServe(":8080", nil)
}
