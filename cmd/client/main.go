package main

import (
	"log"

	"github.com/marcusolsson/tui-go"
)

var logoText = `     
╔═╗╔╦╗╦╔═╗  ╦═╗╔═╗╔═╗  ╦ ╦┌─┐┬─┐┬┌─┌┐ ┌─┐┌┐┌┌─┐┬ ┬
║  ║║║║╚═╗  ╠╦╝╠═╝║    ║║║│ │├┬┘├┴┐├┴┐├┤ ││││  ├─┤
╚═╝╩ ╩╩╚═╝  ╩╚═╩  ╚═╝  ╚╩╝└─┘┴└─┴ ┴└─┘└─┘┘└┘└─┘┴ ┴`

var (
	documentList           *tui.List
	createFolderTextEdit   *tui.TextEdit
	createFolderButton     *tui.Button
	createDocumentTextEdit *tui.TextEdit
	createDocumentButton   *tui.Button
	deleteDocumentButton   *tui.Button
)

func main() {

	window := tui.NewVBox(
		getLogoContainer(),
		getDocumentContainer(),
		getActionContainer(),
	)
	window.SetBorder(true)

	tui.DefaultFocusChain.Set(documentList, createFolderTextEdit, createFolderButton, createDocumentTextEdit, createDocumentButton, deleteDocumentButton)

	ui, err := tui.New(window)
	if err != nil {
		log.Fatal(err)
	}

	ui.SetKeybinding("Esc", func() { ui.Quit() })

	if err := ui.Run(); err != nil {
		log.Fatal(err)
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
	documentList.AddItems("⬆️ ...")
	documentList.SetFocused(true)
	documentWrapper := tui.NewVBox(
		documentList,
	)
	documentWrapper.SetBorder(true)
	documentWrapper.SetTitle(" Documents ")

	propertiesList := tui.NewList()
	propertiesContainer := tui.NewVBox(
		propertiesList,
	)
	propertiesContainer.SetBorder(true)
	propertiesContainer.SetTitle(" Properties ")

	documentContainer := tui.NewHBox(
		documentWrapper,
		propertiesContainer,
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
		tui.NewPadder(1, 1, createFolderContainer),
		tui.NewPadder(1, 0, createDocumentContainer),
		tui.NewPadder(1, 1, deleteDocumentContainer),
	)
	actionContainer.SetBorder(true)
	actionContainer.SetSizePolicy(tui.Minimum, tui.Maximum)
	actionContainer.SetTitle(" Actions ")
	return actionContainer
}
