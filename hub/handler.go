package hub

import (
	"encoding/json"
	"errors"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"log"
	"net/http"
	"wheep-server/storage"
)

func HandleAdd(userId primitive.ObjectID, w http.ResponseWriter, r *http.Request) error {
	var av AddView
	err := json.NewDecoder(r.Body).Decode(&av)
	if err != nil {
		return err
	}
	service := GetService()
	userMap := make(map[primitive.ObjectID]bool)
	users := []primitive.ObjectID{userId}
	for _, v := range av.Users {
		if _, ok := userMap[v]; !ok {
			userMap[v] = true
			users = append(users, v)
		}
	}
	add, err := service.Add(Model{
		Name:  av.Name,
		Image: av.Image,
		Users: users,
	})
	if err != nil {
		return err
	}
	err = json.NewEncoder(w).Encode(add.View())
	if err != nil {
		log.Println(err)
	}
	return nil
}

func HandleGet(userId primitive.ObjectID, w http.ResponseWriter, r *http.Request) error {
	id, err := primitive.ObjectIDFromHex(r.FormValue("id"))
	if err != nil {
		return err
	}
	model, err := GetService().Get(id)
	if err != nil {
		return err
	}
	err = json.NewEncoder(w).Encode(model.View())
	if err != nil {
		log.Println(err)
	}
	return nil
}

func HandleDelete(userId primitive.ObjectID, w http.ResponseWriter, r *http.Request) error {
	id, err := primitive.ObjectIDFromHex(r.FormValue("id"))
	if err != nil {
		return err
	}
	return GetService().Delete(id)
}

func HandleUpdateAvatar(userId primitive.ObjectID, w http.ResponseWriter, r *http.Request) error {
	hubId, err := primitive.ObjectIDFromHex(r.FormValue("hubId"))
	if err != nil {
		return err
	}
	isMember, err := GetService().IsMember(hubId, userId)
	if err != nil {
		return err
	}
	if !isMember {
		return errors.New("you are not a member")
	}
	// 5 MB
	resourceAddress, err := storage.UploadImage(userId, r, 5)
	if err != nil {
		return err
	}
	err = GetService().UpdateAvatar(hubId, resourceAddress)
	if err != nil {
		return err
	}
	_, err = w.Write([]byte("\"" + resourceAddress + "\""))
	return err
}

func HandleRename(userId primitive.ObjectID, w http.ResponseWriter, r *http.Request) error {
	var v View
	err := json.NewDecoder(r.Body).Decode(&v)
	if err != nil {
		return err
	}
	service := GetService()
	err = service.Rename(Model{
		ID:   v.ID,
		Name: v.Name,
	})
	return err
}

func HandleAddUsers(userId primitive.ObjectID, w http.ResponseWriter, r *http.Request) error {
	var users []primitive.ObjectID
	err := json.NewDecoder(r.Body).Decode(&users)
	if err != nil {
		return err
	}
	id, err := primitive.ObjectIDFromHex(r.FormValue("id"))
	if err != nil {
		return err
	}
	service := GetService()
	err = service.AddUsers(id, users)
	return err
}

func HandleRemoveUsers(userId primitive.ObjectID, w http.ResponseWriter, r *http.Request) error {
	var users []primitive.ObjectID
	err := json.NewDecoder(r.Body).Decode(&users)
	if err != nil {
		return err
	}
	id, err := primitive.ObjectIDFromHex(r.FormValue("id"))
	if err != nil {
		return err
	}
	service := GetService()
	err = service.RemoveUsers(id, users)
	return err
}

func HandleFindByUser(userId primitive.ObjectID, w http.ResponseWriter, r *http.Request) error {
	s := r.FormValue("id")
	userId, err := primitive.ObjectIDFromHex(s)
	if err != nil {
		return err
	}
	service := GetService()
	hubs, err := service.FindByUser(userId)
	if err != nil {
		return err
	}
	views := make([]View, len(hubs))
	for i, hub := range hubs {
		views[i] = hub.View()
	}
	return json.NewEncoder(w).Encode(views)
}

func HandleFindMyHubs(userId primitive.ObjectID, w http.ResponseWriter, r *http.Request) error {
	service := GetService()
	hubs, err := service.FindByUser(userId)
	if err != nil {
		return err
	}
	views := make([]View, len(hubs))
	for i, hub := range hubs {
		views[i] = hub.View()
	}
	return json.NewEncoder(w).Encode(views)
}
