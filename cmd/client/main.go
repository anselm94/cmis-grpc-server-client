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
		go loadObject(repository.GetRootFolder().GetId())
	}
}

func loadObject(objectID *cmis.CmisObjectId) {
	ctxt, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	getObjectClient, err := cmisClient.GetObject(ctxt, objectID)
	if err != nil {
		updateStatus(fmt.Sprintf("Error while loading folder %s -> %s", objectID, err))
	} else {
		for {
			cmisObject, err := getObjectClient.Recv()
			if err == io.EOF {
				updateStatus("Server stopped sending updates")
				return
			} else if err != nil {
				updateStatus(fmt.Sprintf("Error while streaming from server -> %s", err))
			}
			updateDocumentList(cmisObject)
			updateStatus(fmt.Sprintf("Streaming the folder \"%s\" from server live", cmisObject.GetId()))
		}
	}
}

func updateStatus(status string) {
	ui.Update(func() {
		statusBar.SetText(status)
	})
}

func updateDocumentList(folderObject *cmis.CmisObject) {
	folder = folderObject
	if folder == nil || folder.Children == nil {
		return
	}
	isRootFolder := proto.Equal(repository.GetRootFolder().GetId(), folder.GetId())
	names := make([]string, len(folder.Children))
	for index, child := range folder.Children {
		var icon string
		var name string
		if child.GetTypeDefinition().GetName() == "cmis:folder" {
			icon = "📂"
		} else {
			icon = "📄"
		}
		for _, property := range child.GetProperties() {
			if property.GetPropertyDefinition().GetName() == "cmis:name" {
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
		go loadObject(folder.Parents[0].GetId())
	} else {
		pos := l.Selected()
		if !isRootFolder {
			pos--
		}
		object := folder.Children[pos]
		if object.TypeDefinition.Name == "cmis:folder" {
			go loadObject(folder.Children[pos].GetId())
		}
	}
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
	deleteDocumentLabel := tui.NewLabel("Delete Selected Document")
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
