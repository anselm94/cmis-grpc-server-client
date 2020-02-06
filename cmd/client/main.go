package main

import (
	"context"
	"docserverclient"
	cmis "docserverclient/proto"
	"fmt"
	"io"
	"log"
	"time"

	"github.com/golang/protobuf/proto"
	"github.com/golang/protobuf/ptypes/empty"
	"github.com/marcusolsson/tui-go"
	"google.golang.org/grpc"
)

const logoText = `     
╔═╗╔╦╗╦╔═╗  ╦═╗╔═╗╔═╗  ╦ ╦┌─┐┬─┐┬┌─┌┐ ┌─┐┌┐┌┌─┐┬ ┬
║  ║║║║╚═╗  ╠╦╝╠═╝║    ║║║│ │├┬┘├┴┐├┴┐├┤ ││││  ├─┤
╚═╝╩ ╩╩╚═╝  ╩╚═╩  ╚═╝  ╚╩╝└─┘┴└─┴ ┴└─┘└─┘┘└┘└─┘┴ ┴`
const navUp = "⬆️ ..."
const typeFolder = "cmis:folder"
const typeDocument = "cmis:document"
const propName = "cmis:name"

var (
	ui                     tui.UI
	statusBar              *tui.StatusBar
	documentList           *tui.List
	propertiesList         *tui.List
	createFolderTextEdit   *tui.TextEdit
	createFolderButton     *tui.Button
	createDocumentTextEdit *tui.TextEdit
	createDocumentButton   *tui.Button
	deleteDocumentButton   *tui.Button

	grpcConnection *grpc.ClientConn
	cmisClient     cmis.CmisServiceClient
	repository     *cmis.Repository
	folder         *cmis.CmisObject
)
var objectIDChannel = make(chan *cmis.CmisObjectId)

func setupUI() {
	statusBar = tui.NewStatusBar("Not connected")
	window := tui.NewVBox(
		getLogoContainer(),
		getDocumentContainer(),
		getActionContainer(),
		statusBar,
	)
	window.SetBorder(true)

	documentList.OnItemActivated(onDocItemSelection)
	createFolderButton.OnActivated(onCreateFolder)
	createDocumentButton.OnActivated(onCreateDocument)
	deleteDocumentButton.OnActivated(onDeleteDocument)

	tui.DefaultFocusChain.Set(documentList, createFolderTextEdit, createFolderButton, createDocumentTextEdit, createDocumentButton, deleteDocumentButton)

	ui, _ = tui.New(window)

	ui.SetKeybinding("Esc", func() { ui.Quit() })

	if err := ui.Run(); err != nil {
		log.Fatalf("Error running UI -> %s", err)
	}
}

func main() {
	config := docserverclient.NewDefaultConfig()

	grpcConnection, err := grpc.Dial(fmt.Sprintf("%s%s", config.AppHost, config.AppPort), grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		log.Fatalf("Connection could not be established with %s:%s -> %s", config.AppHost, config.AppPort, err)
	}
	defer grpcConnection.Close()

	cmisClient = cmis.NewCmisServiceClient(grpcConnection)

	go subscribeObject()
	go loadRepository()

	setupUI()
}

func loadRepository() {
	ctxt, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	repo, err := cmisClient.GetRepository(ctxt, &empty.Empty{})
	if err != nil {
		updateStatus(fmt.Sprintf("Failed to connect -> %s", err))
	} else {
		repository = repo
		updateStatus(fmt.Sprintf("Connected to %s", repository.GetName()))
		objectIDChannel <- repository.GetRootFolder().GetId()
	}
}

func subscribeObject() {
	cmisSubsObjectClient, err := cmisClient.SubscribeObject(context.Background())
	if err != nil {
		updateStatus(fmt.Sprintf("Error establishing a subscription -> %s", err))
	}
	go streamObjectIdsToServer(cmisSubsObjectClient)
	go streamObjectsFromServer(cmisSubsObjectClient)
}

func streamObjectIdsToServer(cmisSubsObjectClient cmis.CmisService_SubscribeObjectClient) {
	for {
		select {
		case cmisObjectID := <-objectIDChannel:
			cmisSubsObjectClient.Send(cmisObjectID)
			updateStatus(fmt.Sprintf("Requested the object for ID - %d", cmisObjectID.Id))
		}
	}
}

func streamObjectsFromServer(cmisSubsObjectClient cmis.CmisService_SubscribeObjectClient) {
	for {
		cmisObject, err := cmisSubsObjectClient.Recv()
		if err == io.EOF {
			updateStatus("Server stopped sending updates")
			return
		} else if err != nil {
			updateStatus(fmt.Sprintf("Error while streaming from server -> %s", err))
			return
		}
		updateDocumentList(cmisObject)
		updateStatus(fmt.Sprintf("Streaming the folder \"%d\" from server live", cmisObject.Id.Id))
	}
}

func createObject(name string, typeStr string) {
	ctxt, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	createRequest := &cmis.CreateObjectReq{
		Name:         name,
		Type:         typeStr,
		ParentId:     folder.GetId(),
		RepositoryId: repository.Id,
	}
	_, err := cmisClient.CreateObject(ctxt, createRequest)
	if err != nil {
		updateStatus(fmt.Sprintf("Failed to create object -> %s", err))
	} else {
		updateStatus(fmt.Sprintf("Created the object"))
	}
}

func deleteObject(objectID *cmis.CmisObjectId) {
	ctxt, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	_, err := cmisClient.DeleteObject(ctxt, objectID)
	if err != nil {
		updateStatus(fmt.Sprintf("Failed to delete object -> %s", err))
	} else {
		updateStatus(fmt.Sprintf("Deleted the object %s", objectID))
	}
}

func updateStatus(status string) {
	ui.Update(func() {
		statusBar.SetText(status)
	})
}

func updateDocumentList(folderObject *cmis.CmisObject) {
	folder = folderObject
	if folder == nil {
		return
	}
	isRootFolder := proto.Equal(repository.GetRootFolder().GetId(), folder.GetId())
	names := make([]string, len(folder.Children))
	for index, child := range folder.Children {
		var icon string
		var name string
		if child.GetTypeDefinition().GetName() == typeFolder {
			icon = "📂"
		} else {
			icon = "📄"
		}
		for _, property := range child.GetProperties() {
			if property.GetPropertyDefinition().GetName() == propName {
				name = property.GetValue()
			}
		}
		names[index] = fmt.Sprintf("%s %s", icon, name)
	}
	ui.Update(func() {
		if documentList.Length() != 0 {
			documentList.RemoveItems()
		}
		if !isRootFolder {
			documentList.AddItems(navUp)
		}
		documentList.AddItems(names...)
	})
}

func onDocItemSelection(l *tui.List) {
	isRootFolder := proto.Equal(repository.GetRootFolder().GetId(), folder.GetId())
	if l.SelectedItem() == navUp {
		objectIDChannel <- folder.Parents[0].GetId()
	} else {
		pos := l.Selected()
		if !isRootFolder {
			pos--
		}
		object := folder.Children[pos]
		if object.TypeDefinition.Name == typeFolder {
			objectIDChannel <- folder.Children[pos].GetId()
		}
	}
}

func onCreateFolder(b *tui.Button) {
	name := createFolderTextEdit.Text()
	if name == "" {
		statusBar.SetText("Enter a folder name!")
		return
	}
	go createObject(name, typeFolder)
}

func onCreateDocument(b *tui.Button) {
	name := createDocumentTextEdit.Text()
	if name == "" {
		statusBar.SetText("Enter a document name!")
		return
	}
	go createObject(name, typeDocument)
}

func onDeleteDocument(b *tui.Button) {
	isRootFolder := proto.Equal(repository.GetRootFolder().GetId(), folder.GetId())
	if documentList.SelectedItem() == navUp {
		statusBar.SetText("You cannot delete your navigator...")
		return
	}
	pos := documentList.Selected()
	if !isRootFolder {
		pos--
	}
	object := folder.Children[pos]
	go deleteObject(object.GetId())
}

func getLogoContainer() *tui.Box {
	logoContainer := tui.NewHBox(
		tui.NewSpacer(),
		tui.NewLabel(logoText),
		tui.NewSpacer(),
	)
	logoContainer.SetSizePolicy(tui.Minimum, tui.Maximum)
	return logoContainer
}

func getDocumentContainer() *tui.Box {
	documentList = tui.NewList()
	documentList.SetFocused(true)
	documentWrapper := tui.NewVBox(
		documentList,
	)
	documentWrapper.SetBorder(true)
	documentWrapper.SetTitle(" Documents ")

	documentContainer := tui.NewHBox(
		documentWrapper,
	)
	return documentContainer
}

func getActionContainer() *tui.Box {
	// Create a Folder
	createFolderLabel := tui.NewLabel("Folder Name : ")
	createFolderTextEdit = tui.NewTextEdit()
	createFolderTextEdit.SetSizePolicy(tui.Minimum, tui.Expanding)
	createFolderButton = tui.NewButton("[ 📂 Create Folder   ]")
	createFolderButton.SetSizePolicy(tui.Minimum, tui.Minimum)
	createFolderContainer := tui.NewHBox(
		createFolderLabel,
		createFolderTextEdit,
		tui.NewSpacer(),
		createFolderButton,
	)

	// Create a Document
	createDocumentLabel := tui.NewLabel("Document Name : ")
	createDocumentTextEdit = tui.NewTextEdit()
	createDocumentTextEdit.SetSizePolicy(tui.Minimum, tui.Expanding)
	createDocumentButton = tui.NewButton("[ 📄 Create Document ]")
	createDocumentButton.SetSizePolicy(tui.Minimum, tui.Minimum)
	createDocumentContainer := tui.NewHBox(
		createDocumentLabel,
		createDocumentTextEdit,
		tui.NewSpacer(),
		createDocumentButton,
	)

	// Delete a Document
	deleteDocumentLabel := tui.NewLabel("Delete Selected : ")
	deleteDocumentButton = tui.NewButton("[ 🚫 Delete Object   ]")
	deleteDocumentButton.SetSizePolicy(tui.Minimum, tui.Minimum)
	deleteDocumentContainer := tui.NewHBox(
		deleteDocumentLabel,
		tui.NewSpacer(),
		deleteDocumentButton,
	)

	actionContainer := tui.NewVBox(
		createFolderContainer,
		createDocumentContainer,
		deleteDocumentContainer,
	)
	actionContainer.SetBorder(true)
	actionContainer.SetSizePolicy(tui.Minimum, tui.Maximum)
	actionContainer.SetTitle(" Actions ")
	return actionContainer
}
