package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"password-manager/controller"
	"password-manager/models"
	"password-manager/security"

	"github.com/therecipe/qt/core"
	"github.com/therecipe/qt/gui"
	"github.com/therecipe/qt/widgets"
)

var tree *widgets.QTreeWidget = nil
var group *widgets.QAction = nil
var sub *widgets.QAction = nil
var add *widgets.QAction = nil
var save *widgets.QAction = nil
var table *widgets.QTableWidget = nil
var fileDB string = ""

func createMenu() *widgets.QMenuBar {
	menu := widgets.NewQMenuBar(nil)
	file := menu.AddMenu2("File")

	newDatabase := widgets.NewQAction(nil)
	newDatabase.SetIcon(gui.NewQIcon5("icons/database.png"))
	newDatabase.SetText("New database")
	newDatabase.ConnectTriggered(func(bool) {
		fmt.Println("New database")
		// controller.CreateDatabase("test")
	})

	openDatabase := widgets.NewQAction(nil)
	openDatabase.SetIcon(gui.NewQIcon5("icons/open.svg"))
	openDatabase.SetText("Open database")
	openDatabase.ConnectTriggered(func(bool) {
		fmt.Println("Open database")
	})

	file.InsertAction(nil, newDatabase)
	file.InsertAction(nil, openDatabase)

	return menu
}

func createToolBar() *widgets.QToolBar {
	tool := widgets.NewQToolBar2(nil)
	tool.SetIconSize(core.NewQSize2(32, 32))
	tool.SetStyleSheet("background-color: #FFFFFF;")
	tool.SetFixedHeight(50)
	database := widgets.NewQAction(nil)
	database.SetIcon(gui.NewQIcon5("icons/database.png"))
	database.SetToolTip("New database")
	database.ConnectTriggered(func(bool) {
		db := newDb()
		if db {
			file := saveFile()
			if file != "" && !controller.CheckFileExist(file) {
				name := filepath.Base(file)
				name2 := strings.TrimSuffix(name, filepath.Ext(name))
				password := createPassword(file)
				if password != "" {
					init := controller.InitDB(file)
					if init != nil {
						log.Println(init)
						showError("Failed to init database!")
						return
					}
					create := controller.CreateDatabaseAndSecretGroupIfNotExist(file, name2)
					if create != nil {
						log.Println(create)
						showError("Failed to create database!")
						return
					}
					databases, err := controller.GetAllDatabases(file)
					if err != nil {
						log.Println(err)
						showError("Failed to get data!")
						return
					}
					tree.Clear()
					table.ClearContents()
					table.SetRowCount(0)
					for _, database := range databases {
						parent := widgets.NewQTreeWidgetItem2([]string{database.Name}, 0)
						parent.SetIcon(0, gui.NewQIcon5("icons/sub.svg"))
						tree.AddTopLevelItem(parent)
						for _, group := range database.SecretGroups {
							child := widgets.NewQTreeWidgetItem2([]string{group.Name}, 0)
							child.SetIcon(0, gui.NewQIcon5("icons/group.svg"))
							parent.AddChild(child)
						}
						parent.SetExpanded(true)
					}
					group.SetEnabled(true)
					add.SetEnabled(true)
					save.SetEnabled(true)
					sub.SetEnabled(true)
					fileDB = file
				}
			} else {
				log.Println("File exists or cancelled")
			}
		} else {
			log.Println("Cancel")
		}
	})

	open := widgets.NewQAction(nil)
	open.SetIcon(gui.NewQIcon5("icons/open.svg"))
	open.SetToolTip("Open database")
	open.ConnectTriggered(func(bool) {
		file := loadFile()
		if file != "" {
			databases, err := controller.GetAllDatabases(file)
			if err != nil {
				log.Println(err)
				showError("Failed to get data!")
				return
			}
			tree.Clear()
			table.ClearContents()
			table.SetRowCount(0)
			for _, database := range databases {
				parent := widgets.NewQTreeWidgetItem2([]string{database.Name}, 0)
				parent.SetIcon(0, gui.NewQIcon5("icons/sub.svg"))
				tree.AddTopLevelItem(parent)
				for i, group := range database.SecretGroups {
					child := widgets.NewQTreeWidgetItem2([]string{group.Name}, 0)
					child.SetIcon(0, gui.NewQIcon5("icons/group.svg"))
					parent.AddChild(child)
					if i == 0 {
						tree.SetCurrentItem(child)
						for _, secret := range group.Secrets {
							row := table.RowCount()
							table.InsertRow(row)
							title := widgets.NewQTableWidgetItem2(secret.Title, 0)
							title.SetIcon(gui.NewQIcon5("icons/key.png"))
							table.SetItem(row, 0, widgets.NewQTableWidgetItem2(fmt.Sprint(secret.ID), 0))
							table.SetItem(row, 1, title)
							table.SetItem(row, 2, widgets.NewQTableWidgetItem2(secret.Username, 0))
							table.SetItem(row, 3, widgets.NewQTableWidgetItem2(secret.Password, 0))
							table.SetItem(row, 4, widgets.NewQTableWidgetItem2(secret.URL, 0))
							table.SetItem(row, 5, widgets.NewQTableWidgetItem2(secret.Description, 0))
						}
					}
				}
				parent.SetExpanded(true)
			}
			group.SetEnabled(true)
			add.SetEnabled(true)
			save.SetEnabled(true)
			sub.SetEnabled(true)
			fileDB = file
		} else {
			log.Println("Cancel")
		}
	})

	sub = widgets.NewQAction(nil)
	sub.SetIcon(gui.NewQIcon5("icons/sub.svg"))
	sub.SetToolTip("Create sub database")
	sub.ConnectTriggered(func(bool) {
		log.Println("Create sub database")
		db := getSubDatabaseName("")
		if db != "" {
			database, err := controller.CreateSubDatabase(fileDB, db)
			if err != nil {
				log.Println(err)
				showError("Failed to add sub database!")
				return
			} else {
				db, err := controller.GetDatabase(fileDB, database.Name)
				if err != nil {
					log.Println(err)
					showError("Failed to get data!")
					return
				}
				parent := widgets.NewQTreeWidgetItem2([]string{db.Name}, 0)
				parent.SetIcon(0, gui.NewQIcon5("icons/sub.svg"))
				tree.AddTopLevelItem(parent)
				for _, group := range db.SecretGroups {
					child := widgets.NewQTreeWidgetItem2([]string{group.Name}, 0)
					child.SetIcon(0, gui.NewQIcon5("icons/group.svg"))
					parent.AddChild(child)
				}
				parent.SetExpanded(true)
			}
		}
	})
	sub.SetEnabled(false)

	group = widgets.NewQAction(nil)
	group.SetIcon(gui.NewQIcon5("icons/group.svg"))
	group.SetToolTip("Add new secret group")
	group.ConnectTriggered(func(bool) {
		grp := getSecretGroup("")
		if grp != "" {
			if tree.CurrentItem().Parent().Text(0) == "" {
				grp, err := controller.CreateSecretGroup(fileDB, tree.CurrentItem().Text(0), grp)
				if err != nil {
					log.Println(err)
					showError("Failed to add secret group!")
					return
				} else {
					g := widgets.NewQTreeWidgetItem2([]string{grp.Name}, 0)
					g.SetIcon(0, gui.NewQIcon5("icons/group.svg"))
					tree.CurrentItem().AddChild(g)
				}
			} else {
				grp, err := controller.CreateSecretGroup(fileDB, tree.CurrentItem().Parent().Text(0), grp)
				if err != nil {
					log.Println(err)
					showError("Failed to add secret group!")
					return
				} else {
					g := widgets.NewQTreeWidgetItem2([]string{grp.Name}, 0)
					g.SetIcon(0, gui.NewQIcon5("icons/group.svg"))
					tree.CurrentItem().Parent().AddChild(g)
				}
			}
		}
	})
	group.SetEnabled(false)

	add = widgets.NewQAction(nil)
	add.SetIcon(gui.NewQIcon5("icons/key.png"))
	add.SetToolTip("Add new secret")
	add.ConnectTriggered(func(bool) {
		secret := getSecret(models.Secret{})
		if secret == (models.Secret{}) {
			return
		}
		if tree.CurrentItem().Parent().Text(0) == "" {
			if tree.CurrentItem().Child(0).Text(0) != "" {
				sct, err := controller.CreateSecret(fileDB, tree.CurrentItem().Text(0), tree.CurrentItem().Child(0).Text(0), secret)
				if err != nil {
					log.Println(err)
					showError("Failed to add secret!")
					return
				} else {
					row := table.RowCount()
					table.InsertRow(row)
					title := widgets.NewQTableWidgetItem2(sct.Title, 0)
					title.SetIcon(gui.NewQIcon5("icons/key.png"))
					table.SetItem(row, 0, widgets.NewQTableWidgetItem2(fmt.Sprint(sct.ID), 0))
					table.SetItem(row, 1, title)
					table.SetItem(row, 2, widgets.NewQTableWidgetItem2(sct.Username, 0))
					table.SetItem(row, 3, widgets.NewQTableWidgetItem2(sct.Password, 0))
					table.SetItem(row, 4, widgets.NewQTableWidgetItem2(sct.URL, 0))
					table.SetItem(row, 5, widgets.NewQTableWidgetItem2(sct.Description, 0))
				}
			}
		} else {
			sct, err := controller.CreateSecret(fileDB, tree.CurrentItem().Parent().Text(0), tree.CurrentItem().Text(0), secret)
			if err != nil {
				log.Println(err)
				showError("Failed to add secret!")
				return
			} else {
				row := table.RowCount()
				table.InsertRow(row)
				title := widgets.NewQTableWidgetItem2(sct.Title, 0)
				title.SetIcon(gui.NewQIcon5("icons/key.png"))
				table.SetItem(row, 0, widgets.NewQTableWidgetItem2(fmt.Sprint(sct.ID), 0))
				table.SetItem(row, 1, title)
				table.SetItem(row, 2, widgets.NewQTableWidgetItem2(sct.Username, 0))
				table.SetItem(row, 3, widgets.NewQTableWidgetItem2(sct.Password, 0))
				table.SetItem(row, 4, widgets.NewQTableWidgetItem2(sct.URL, 0))
				table.SetItem(row, 5, widgets.NewQTableWidgetItem2(sct.Description, 0))
			}
		}
	})
	add.SetEnabled(false)

	save = widgets.NewQAction(nil)
	save.SetIcon(gui.NewQIcon5("icons/save.svg"))
	save.SetToolTip("Save")
	save.ConnectTriggered(func(bool) {
		log.Println("Save")
	})
	save.SetEnabled(false)

	line := widgets.NewQFrame(nil, 0)
	line.SetFrameShape(widgets.QFrame__VLine)
	line.SetFrameShadow(widgets.QFrame__Sunken)

	line2 := widgets.NewQFrame(nil, 0)
	line2.SetFrameShape(widgets.QFrame__VLine)
	line2.SetFrameShadow(widgets.QFrame__Sunken)

	tool.InsertAction(nil, database)
	tool.InsertAction(nil, open)
	tool.AddWidget(line)
	tool.InsertAction(nil, save)
	tool.AddWidget(line2)
	tool.InsertAction(nil, sub)
	tool.InsertAction(nil, group)
	tool.InsertAction(nil, add)

	return tool
}

func createSideMenu() *widgets.QWidget {
	widget := widgets.NewQWidget(nil, 0)
	widget.SetFixedWidth(100)
	widget.SetMinimumWidth(50)
	widget.SetMaximumWidth(200)

	tree = widgets.NewQTreeWidget(nil)
	tree.SetHeaderHidden(true)
	tree.SetColumnCount(1)
	tree.SetAnimated(true)
	tree.SetUniformRowHeights(true)
	tree.SetItemsExpandable(true)
	tree.SetHorizontalScrollBarPolicy(1)
	tree.SetVerticalScrollBarPolicy(1)
	tree.SetAutoScroll(true)
	tree.SetAutoScrollMargin(10)

	tree.ConnectItemClicked(func(item *widgets.QTreeWidgetItem, column int) {
		table.ClearContents()
		table.SetRowCount(0)
		if item.Parent().Text(0) == "" {
			if item.Child(0).Text(0) != "" {
				secrets, err := controller.GetSecrets(fileDB, item.Text(0), item.Child(0).Text(0))
				if err != nil {
					log.Println(err)
					showError("Failed to get data!")
					return
				}
				for _, secret := range secrets {
					row := table.RowCount()
					table.InsertRow(row)
					title := widgets.NewQTableWidgetItem2(secret.Title, 0)
					title.SetIcon(gui.NewQIcon5("icons/key.png"))
					table.SetItem(row, 0, widgets.NewQTableWidgetItem2(fmt.Sprint(secret.ID), 0))
					table.SetItem(row, 1, title)
					table.SetItem(row, 2, widgets.NewQTableWidgetItem2(secret.Username, 0))
					table.SetItem(row, 3, widgets.NewQTableWidgetItem2(secret.Password, 0))
					table.SetItem(row, 4, widgets.NewQTableWidgetItem2(secret.URL, 0))
					table.SetItem(row, 5, widgets.NewQTableWidgetItem2(secret.Description, 0))
				}
			}
		} else {
			secrets, err := controller.GetSecrets(fileDB, item.Parent().Text(0), item.Text(0))
			if err != nil {
				log.Println(err)
				showError("Failed to get data!")
				return
			}
			for _, secret := range secrets {
				row := table.RowCount()
				table.InsertRow(row)
				title := widgets.NewQTableWidgetItem2(secret.Title, 0)
				title.SetIcon(gui.NewQIcon5("icons/key.png"))
				table.SetItem(row, 0, widgets.NewQTableWidgetItem2(fmt.Sprint(secret.ID), 0))
				table.SetItem(row, 1, title)
				table.SetItem(row, 2, widgets.NewQTableWidgetItem2(secret.Username, 0))
				table.SetItem(row, 3, widgets.NewQTableWidgetItem2(secret.Password, 0))
				table.SetItem(row, 4, widgets.NewQTableWidgetItem2(secret.URL, 0))
				table.SetItem(row, 5, widgets.NewQTableWidgetItem2(secret.Description, 0))
			}
		}
	})

	menu := widgets.NewQMenu(nil)

	edit := menu.AddAction("Edit")
	edit.SetIcon(gui.NewQIcon5("icons/edit.svg"))
	edit.ConnectTriggered(func(bool) {
		if tree.CurrentItem().Parent().Text(0) == "" {
			name := getSubDatabaseName(tree.CurrentItem().Text(0))
			if name != "" {
				d, err := controller.UpdateDatabase(fileDB, tree.CurrentItem().Text(0), name)
				if err != nil {
					log.Println(err)
					showError("Failed to update database!")
					return
				} else {
					tree.CurrentItem().SetText(0, d.Name)
				}
			}
		} else {
			group := getSecretGroup(tree.CurrentItem().Text(0))
			if group != "" {
				g, err := controller.UpdateSecretGroup(fileDB, tree.CurrentItem().Parent().Text(0), tree.CurrentItem().Text(0), group)
				if err != nil {
					log.Println(err)
					showError("Failed to update secret group!")
					return
				} else {
					tree.CurrentItem().SetText(0, g.Name)
				}
			}
		}
	})

	delete := menu.AddAction("Delete")
	delete.SetIcon(gui.NewQIcon5("icons/delete.svg"))
	delete.ConnectTriggered(func(bool) {
		if tree.CurrentItem().Parent().Text(0) == "" {
			if tree.CurrentItem().Child(0).Text(0) != "" {
				showError("Delete all sub databases first!")
			} else {
				err := controller.DeleteDatabase(fileDB, tree.CurrentItem().Text(0))
				if err != nil {
					log.Println(err)
					showError("Failed to delete database!")
					return
				} else {
					tree.Clear()
					table.ClearContents()
					table.SetRowCount(0)
					group.SetEnabled(false)
					add.SetEnabled(false)
					save.SetEnabled(false)
					sub.SetEnabled(false)
				}
			}
		} else {
			group, err := controller.GetSecretGroup(fileDB, tree.CurrentItem().Parent().Text(0), tree.CurrentItem().Text(0))
			if err != nil {
				log.Println(err)
				showError("Failed to delete secret group!")
				return
			}
			err = controller.DeleteSecretGroup(fileDB, tree.CurrentItem().Parent().Text(0), group.Name)
			if err != nil {
				log.Println(err)
				showError("Failed to delete secret group!")
				return
			} else {
				tree.CurrentItem().Parent().RemoveChild(tree.CurrentItem())
			}
		}
	})

	tree.SetContextMenuPolicy(core.Qt__CustomContextMenu)

	tree.ConnectCustomContextMenuRequested(func(pos *core.QPoint) {
		menu.Exec2(tree.MapToGlobal(pos), nil)
	})

	layout := widgets.NewQHBoxLayout2(widget)
	layout.SetContentsMargins(0, 0, 0, 0)
	layout.SetSpacing(0)
	layout.AddWidget(tree, 0, 0)

	return widget
}

func createMain() *widgets.QWidget {
	widget := widgets.NewQWidget(nil, 0)
	widget.SetStyleSheet("background-color: #FFFFFF;")

	table = widgets.NewQTableWidget(nil)
	table.SetColumnCount(6)
	table.SetRowCount(0)
	table.SetHorizontalHeaderLabels([]string{"ID", "Title", "Username", "Password", "URL", "Description"})
	table.SetEditTriggers(widgets.QAbstractItemView__NoEditTriggers)
	table.SetSelectionBehavior(widgets.QAbstractItemView__SelectRows)
	table.SetSelectionMode(widgets.QAbstractItemView__SingleSelection)
	table.SetShowGrid(true)
	table.SetHorizontalScrollBarPolicy(1)
	table.SetVerticalScrollBarPolicy(1)
	table.SetAutoScroll(true)
	table.SetAutoScrollMargin(10)
	table.SetColumnHidden(0, true)
	table.VerticalHeader().SetVisible(false)

	table.ConnectCellDoubleClicked(func(row int, column int) {
		id := table.Item(row, 0).Text()
		integer, err := strconv.Atoi(id)
		if err != nil {
			log.Println(err)
			showError("Failed to update secret!")
			return
		}
		if tree.CurrentItem().Parent().Text(0) == "" {
			if tree.CurrentItem().Child(0).Text(0) != "" {
				s, err := controller.GetSecret(fileDB, tree.CurrentItem().Text(0), tree.CurrentItem().Child(0).Text(0), integer)
				if err != nil {
					log.Println(err)
					showError("Failed to update secret!")
					return
				}
				secret := getSecret(s)
				if secret == (models.Secret{}) {
					return
				}
				sct, err := controller.UpdateSecret(fileDB, tree.CurrentItem().Text(0), tree.CurrentItem().Child(0).Text(0), integer, secret)
				if err != nil {
					log.Println(err)
					showError("Failed to update secret!")
					return
				} else {
					title := widgets.NewQTableWidgetItem2(sct.Title, 0)
					title.SetIcon(gui.NewQIcon5("icons/key.png"))
					table.SetItem(row, 0, widgets.NewQTableWidgetItem2(fmt.Sprint(sct.ID), 0))
					table.SetItem(row, 1, title)
					table.SetItem(row, 2, widgets.NewQTableWidgetItem2(sct.Username, 0))
					table.SetItem(row, 3, widgets.NewQTableWidgetItem2(sct.Password, 0))
					table.SetItem(row, 4, widgets.NewQTableWidgetItem2(sct.URL, 0))
					table.SetItem(row, 5, widgets.NewQTableWidgetItem2(sct.Description, 0))
				}
			}
		} else {
			s, err := controller.GetSecret(fileDB, tree.CurrentItem().Parent().Text(0), tree.CurrentItem().Text(0), integer)
			if err != nil {
				log.Println(err)
				showError("Failed to update secret!")
				return
			}
			secret := getSecret(s)
			if secret == (models.Secret{}) {
				return
			}
			sct, err := controller.UpdateSecret(fileDB, tree.CurrentItem().Parent().Text(0), tree.CurrentItem().Text(0), integer, secret)
			if err != nil {
				log.Println(err)
				showError("Failed to update secret!")
				return
			} else {
				title := widgets.NewQTableWidgetItem2(sct.Title, 0)
				title.SetIcon(gui.NewQIcon5("icons/key.png"))
				table.SetItem(row, 0, widgets.NewQTableWidgetItem2(fmt.Sprint(sct.ID), 0))
				table.SetItem(row, 1, title)
				table.SetItem(row, 2, widgets.NewQTableWidgetItem2(sct.Username, 0))
				table.SetItem(row, 3, widgets.NewQTableWidgetItem2(sct.Password, 0))
				table.SetItem(row, 4, widgets.NewQTableWidgetItem2(sct.URL, 0))
				table.SetItem(row, 5, widgets.NewQTableWidgetItem2(sct.Description, 0))
			}
		}
	})

	menu := widgets.NewQMenu(nil)

	edit := menu.AddAction("Edit")
	edit.SetIcon(gui.NewQIcon5("icons/edit.svg"))
	edit.ConnectTriggered(func(bool) {
		row := table.CurrentRow()
		id := table.Item(row, 0).Text()
		integer, err := strconv.Atoi(id)
		if err != nil {
			log.Println(err)
			showError("Failed to update secret!")
			return
		}
		if tree.CurrentItem().Parent().Text(0) == "" {
			if tree.CurrentItem().Child(0).Text(0) != "" {
				s, err := controller.GetSecret(fileDB, tree.CurrentItem().Text(0), tree.CurrentItem().Child(0).Text(0), integer)
				if err != nil {
					log.Println(err)
					showError("Failed to update secret!")
					return
				}
				secret := getSecret(s)
				if secret == (models.Secret{}) {
					return
				}
				sct, err := controller.UpdateSecret(fileDB, tree.CurrentItem().Text(0), tree.CurrentItem().Child(0).Text(0), integer, secret)
				if err != nil {
					log.Println(err)
					showError("Failed to update secret!")
					return
				} else {
					title := widgets.NewQTableWidgetItem2(sct.Title, 0)
					title.SetIcon(gui.NewQIcon5("icons/key.png"))
					table.SetItem(row, 0, widgets.NewQTableWidgetItem2(fmt.Sprint(sct.ID), 0))
					table.SetItem(row, 1, title)
					table.SetItem(row, 2, widgets.NewQTableWidgetItem2(sct.Username, 0))
					table.SetItem(row, 3, widgets.NewQTableWidgetItem2(sct.Password, 0))
					table.SetItem(row, 4, widgets.NewQTableWidgetItem2(sct.URL, 0))
					table.SetItem(row, 5, widgets.NewQTableWidgetItem2(sct.Description, 0))
				}
			}
		} else {
			s, err := controller.GetSecret(fileDB, tree.CurrentItem().Parent().Text(0), tree.CurrentItem().Text(0), integer)
			if err != nil {
				log.Println(err)
				showError("Failed to update secret!")
				return
			}
			secret := getSecret(s)
			if secret == (models.Secret{}) {
				return
			}
			sct, err := controller.UpdateSecret(fileDB, tree.CurrentItem().Parent().Text(0), tree.CurrentItem().Text(0), integer, secret)
			if err != nil {
				log.Println(err)
				showError("Failed to update secret!")
				return
			} else {
				title := widgets.NewQTableWidgetItem2(sct.Title, 0)
				title.SetIcon(gui.NewQIcon5("icons/key.png"))
				table.SetItem(row, 0, widgets.NewQTableWidgetItem2(fmt.Sprint(sct.ID), 0))
				table.SetItem(row, 1, title)
				table.SetItem(row, 2, widgets.NewQTableWidgetItem2(sct.Username, 0))
				table.SetItem(row, 3, widgets.NewQTableWidgetItem2(sct.Password, 0))
				table.SetItem(row, 4, widgets.NewQTableWidgetItem2(sct.URL, 0))
				table.SetItem(row, 5, widgets.NewQTableWidgetItem2(sct.Description, 0))
			}
		}
	})

	delete := menu.AddAction("Delete")
	delete.SetIcon(gui.NewQIcon5("icons/delete.svg"))
	delete.ConnectTriggered(func(bool) {
		row := table.CurrentRow()
		id := table.Item(row, 0).Text()
		integer, err := strconv.Atoi(id)
		if err != nil {
			log.Println(err)
			showError("Failed to delete secret!")
			return
		}
		if tree.CurrentItem().Parent().Text(0) == "" {
			if tree.CurrentItem().Child(0).Text(0) != "" {
				err := controller.DeleteSecret(fileDB, tree.CurrentItem().Text(0), tree.CurrentItem().Child(0).Text(0), integer)
				if err != nil {
					log.Println(err)
					showError("Failed to delete secret!")
					return
				}
				table.RemoveRow(row)
			}
		} else {
			err := controller.DeleteSecret(fileDB, tree.CurrentItem().Parent().Text(0), tree.CurrentItem().Text(0), integer)
			if err != nil {
				log.Println(err)
				showError("Failed to delete secret!")
				return
			}
			table.RemoveRow(row)
		}
	})

	table.SetContextMenuPolicy(core.Qt__CustomContextMenu)

	table.ConnectCustomContextMenuRequested(func(pos *core.QPoint) {
		menu.Exec2(table.MapToGlobal(pos), nil)
	})

	layout := widgets.NewQHBoxLayout2(widget)
	layout.SetContentsMargins(0, 0, 0, 0)
	layout.SetSpacing(0)
	layout.AddWidget(table, 0, 0)

	return widget
}

func createPassword(file string) string {
	dialog := widgets.NewQDialog(nil, 0)
	dialog.SetWindowTitle("Create master key")

	layout := widgets.NewQVBoxLayout2(dialog)

	horizontalLayout := widgets.NewQHBoxLayout2(nil)

	formLayout := widgets.NewQFormLayout(nil)

	label := widgets.NewQLabel(nil, 0)
	label.SetText(fmt.Sprintf("Database: %s\n\nSpecify a new master key, which will be used to encrypt the database.\n\nRemember the master key that you enter, \nif you lose it you will not be able to open the database.", file))

	passwordField := widgets.NewQLineEdit(nil)
	repeatField := widgets.NewQLineEdit(nil)

	passwordField.SetEchoMode(2)

	repeatField.SetEchoMode(2)

	passwordField.ConnectTextChanged(func(_ string) {
		if passwordField.Text() != repeatField.Text() {
			passwordField.SetStyleSheet("background-color: red")
			repeatField.SetStyleSheet("background-color: red")
		} else {
			passwordField.SetStyleSheet("background-color: green")
			repeatField.SetStyleSheet("background-color: green")
		}
	})

	repeatField.ConnectTextChanged(func(_ string) {
		if passwordField.Text() != repeatField.Text() {
			passwordField.SetStyleSheet("background-color: red")
			repeatField.SetStyleSheet("background-color: red")
		} else {
			passwordField.SetStyleSheet("background-color: green")
			repeatField.SetStyleSheet("background-color: green")
		}
	})

	formLayout.AddRow5(label)
	formLayout.AddRow3("Master password:", passwordField)
	formLayout.AddRow3("Repeat password:", repeatField)

	formLayout2 := widgets.NewQFormLayout(nil)

	sh := gui.NewQIcon5("icons/show.svg")
	show := widgets.NewQPushButton3(sh, "", nil)
	show.SetStyleSheet("border-width: 0px;")
	show.ConnectClicked(func(bool) {
		if passwordField.EchoMode() == 2 {
			passwordField.SetEchoMode(0)
			repeatField.SetEchoMode(0)
			sh = gui.NewQIcon5("icons/dontshow.png")
			show.SetIcon(sh)
		} else {
			passwordField.SetEchoMode(2)
			repeatField.SetEchoMode(2)
			sh = gui.NewQIcon5("icons/show.svg")
			show.SetIcon(sh)
		}
	})

	button := widgets.NewQPushButton3(gui.NewQIcon5("icons/refresh.svg"), "", nil)
	button.SetStyleSheet("border-width: 0px;")

	password := ""

	button.ConnectClicked(func(bool) {
		password = security.GenerateStrongPassword(20)
		passwordField.SetText(password)
		repeatField.SetText(password)
		passwordField.SetStyleSheet("background-color: green")
		repeatField.SetStyleSheet("background-color: green")
	})

	formLayout2.AddRow5(widgets.NewQWidget(nil, 0))
	formLayout2.AddRow5(widgets.NewQWidget(nil, 0))
	formLayout2.AddRow5(widgets.NewQWidget(nil, 0))
	formLayout2.AddRow5(widgets.NewQWidget(nil, 0))
	formLayout2.AddRow5(widgets.NewQWidget(nil, 0))
	formLayout2.AddRow5(widgets.NewQWidget(nil, 0))
	formLayout2.AddRow5(widgets.NewQWidget(nil, 0))
	formLayout2.AddRow5(widgets.NewQWidget(nil, 0))
	formLayout2.AddRow5(widgets.NewQWidget(nil, 0))
	formLayout2.AddRow5(widgets.NewQWidget(nil, 0))
	formLayout2.AddRow5(widgets.NewQWidget(nil, 0))
	formLayout2.AddRow5(widgets.NewQWidget(nil, 0))
	formLayout2.AddRow5(widgets.NewQWidget(nil, 0))
	formLayout2.AddRow5(widgets.NewQWidget(nil, 0))
	formLayout2.AddRow5(widgets.NewQWidget(nil, 0))
	formLayout2.AddRow5(show)
	formLayout2.AddRow5(button)

	horizontalLayout.AddLayout(formLayout, 0)
	horizontalLayout.AddLayout(formLayout2, 0)

	layout.AddLayout(horizontalLayout, 0)

	buttons := widgets.NewQDialogButtonBox(nil)
	buttons.SetOrientation(core.Qt__Horizontal)
	buttons.SetStandardButtons(widgets.QDialogButtonBox__Ok | widgets.QDialogButtonBox__Cancel)
	buttons.ConnectAccepted(func() {
		if passwordField.Text() == repeatField.Text() {
			dialog.Accept()
		} else {
			showError("Passwords do not match!")
		}
	})
	buttons.ConnectRejected(func() {
		dialog.Reject()
	})
	layout.AddWidget(buttons, 0, core.Qt__AlignRight)

	dialog.SetModal(true)
	dialog.Show()

	if dialog.Exec() == int(widgets.QDialog__Accepted) {
		return passwordField.Text()
	}
	return ""
}

func getPassword(file string) string {
	dialog := widgets.NewQInputDialog(nil, 0)
	dialog.SetWindowTitle("Enter master key")
	dialog.SetLabelText(fmt.Sprintf("Database: %s\n\nEnter master key to open database.", file))
	dialog.SetOkButtonText("Ok")
	dialog.SetCancelButtonText("Cancel")
	dialog.SetTextEchoMode(2)
	dialog.SetInputMode(widgets.QInputDialog__TextInput)

	checkbox := widgets.NewQCheckBox(dialog)
	checkbox.SetText("Show Password")
	checkbox.SetChecked(false)
	checkbox.ConnectStateChanged(func(state int) {
		if state == 2 {
			dialog.SetTextEchoMode(0)
		} else {
			dialog.SetTextEchoMode(2)
		}
	})
	dialog.Layout().AddWidget(checkbox)

	dialog.SetModal(true)
	dialog.Show()
	if dialog.Exec() == 1 {
		return dialog.TextValue()
	}
	return ""
}

func newDb() bool {
	dialog := widgets.NewQMessageBox(nil)
	dialog.SetWindowTitle("New database")
	dialog.SetText("Your data will be stored in a regular file.\n\nAfter clicking OK, you will be asked to choose a location where the file should be saved.")
	dialog.SetIcon(widgets.QMessageBox__Information)
	dialog.SetStandardButtons(widgets.QMessageBox__Ok | widgets.QMessageBox__Cancel)
	dialog.SetDefaultButton2(widgets.QMessageBox__Ok)
	dialog.SetEscapeButton2(widgets.QMessageBox__Cancel)
	dialog.SetModal(true)
	dialog.Show()
	return dialog.Exec() == int(widgets.QMessageBox__Ok)
}

func getSubDatabaseName(name string) string {
	dialog := widgets.NewQInputDialog(nil, 0)
	dialog.SetWindowTitle("Create sub database")
	dialog.SetLabelText("Specify a new name for the sub database.")
	dialog.SetOkButtonText("Ok")
	dialog.SetCancelButtonText("Cancel")
	dialog.SetTextEchoMode(0)
	dialog.SetInputMode(widgets.QInputDialog__TextInput)
	dialog.SetModal(true)

	if name != "" {
		dialog.SetWindowTitle("Edit sub database")
		dialog.SetLabelText("Specify a new name for the sub database.")
		dialog.SetTextValue(name)
	}

	dialog.Show()
	if dialog.Exec() == 1 {
		return dialog.TextValue()
	}
	return ""
}

func getSecretGroup(group string) string {
	dialog := widgets.NewQInputDialog(nil, 0)
	dialog.SetWindowTitle("Create secret group")
	dialog.SetLabelText("Specify a name for the secret group.")
	dialog.SetOkButtonText("Ok")
	dialog.SetCancelButtonText("Cancel")
	dialog.SetTextEchoMode(0)
	dialog.SetInputMode(widgets.QInputDialog__TextInput)
	dialog.SetModal(true)

	if group != "" {
		dialog.SetWindowTitle("Edit secret group")
		dialog.SetLabelText("Specify a new name for the secret group.")
		dialog.SetTextValue(group)
	}

	dialog.Show()
	if dialog.Exec() == 1 {
		return dialog.TextValue()
	}
	return ""
}

func getSecret(secret models.Secret) models.Secret {
	dialog := widgets.NewQDialog(nil, 0)
	dialog.SetWindowTitle("Create secret")

	layout := widgets.NewQVBoxLayout2(dialog)

	horizontalLayout := widgets.NewQHBoxLayout2(nil)

	formLayout := widgets.NewQFormLayout(nil)

	password := security.GenerateStrongPassword(20)

	titleField := widgets.NewQLineEdit(nil)
	usernameField := widgets.NewQLineEdit(nil)
	passwordField := widgets.NewQLineEdit(nil)
	repeatField := widgets.NewQLineEdit(nil)
	urlField := widgets.NewQLineEdit(nil)
	descriptionField := widgets.NewQTextEdit(nil)

	usernameField.SetStyleSheet("background-color: red")

	passwordField.SetText(password)
	passwordField.SetEchoMode(2)
	passwordField.SetStyleSheet("background-color: green")

	repeatField.SetText(password)
	repeatField.SetEchoMode(2)
	repeatField.SetStyleSheet("background-color: green")

	usernameField.ConnectTextChanged(func(_ string) {
		if usernameField.Text() != "" {
			usernameField.SetStyleSheet("background-color: green")
		} else {
			usernameField.SetStyleSheet("background-color: red")
		}
	})

	passwordField.ConnectTextChanged(func(_ string) {
		if passwordField.Text() != repeatField.Text() {
			passwordField.SetStyleSheet("background-color: red")
			repeatField.SetStyleSheet("background-color: red")
		} else {
			passwordField.SetStyleSheet("background-color: green")
			repeatField.SetStyleSheet("background-color: green")
		}
	})

	repeatField.ConnectTextChanged(func(_ string) {
		if passwordField.Text() != repeatField.Text() {
			passwordField.SetStyleSheet("background-color: red")
			repeatField.SetStyleSheet("background-color: red")
		} else {
			passwordField.SetStyleSheet("background-color: green")
			repeatField.SetStyleSheet("background-color: green")
		}
	})

	if secret != (models.Secret{}) {
		dialog.SetWindowTitle("Edit secret")
		titleField.SetText(secret.Title)
		usernameField.SetText(secret.Username)
		passwordField.SetText(secret.Password)
		repeatField.SetText(secret.Password)
		urlField.SetText(secret.URL)
		descriptionField.SetText(secret.Description)
		if usernameField.Text() != "" {
			usernameField.SetStyleSheet("background-color: green")
		} else {
			usernameField.SetStyleSheet("background-color: red")
		}
	}

	formLayout.AddRow3("Title:", titleField)
	formLayout.AddRow3("Username:", usernameField)
	formLayout.AddRow3("Password:", passwordField)
	formLayout.AddRow3("Repeat password:", repeatField)
	formLayout.AddRow3("URL:", urlField)
	formLayout.AddRow3("Description:", descriptionField)

	formLayout2 := widgets.NewQFormLayout(nil)

	sh := gui.NewQIcon5("icons/show.svg")
	show := widgets.NewQPushButton3(sh, "", nil)
	show.SetStyleSheet("border-width: 0px;")
	show.ConnectClicked(func(bool) {
		if passwordField.EchoMode() == 2 {
			passwordField.SetEchoMode(0)
			repeatField.SetEchoMode(0)
			sh = gui.NewQIcon5("icons/dontshow.png")
			show.SetIcon(sh)
		} else {
			passwordField.SetEchoMode(2)
			repeatField.SetEchoMode(2)
			sh = gui.NewQIcon5("icons/show.svg")
			show.SetIcon(sh)
		}
	})

	button := widgets.NewQPushButton3(gui.NewQIcon5("icons/refresh.svg"), "", nil)
	button.SetStyleSheet("border-width: 0px;")

	button.ConnectClicked(func(bool) {
		password = security.GenerateStrongPassword(20)
		passwordField.SetText(password)
		repeatField.SetText(password)
		passwordField.SetStyleSheet("background-color: green")
		repeatField.SetStyleSheet("background-color: green")
	})

	formLayout2.AddRow5(widgets.NewQWidget(nil, 0))
	formLayout2.AddRow5(widgets.NewQWidget(nil, 0))
	formLayout2.AddRow5(widgets.NewQWidget(nil, 0))
	formLayout2.AddRow5(widgets.NewQWidget(nil, 0))
	formLayout2.AddRow5(widgets.NewQWidget(nil, 0))
	formLayout2.AddRow5(widgets.NewQWidget(nil, 0))
	formLayout2.AddRow5(widgets.NewQWidget(nil, 0))
	formLayout2.AddRow5(widgets.NewQWidget(nil, 0))
	formLayout2.AddRow5(widgets.NewQWidget(nil, 0))
	formLayout2.AddRow5(widgets.NewQWidget(nil, 0))
	formLayout2.AddRow5(show)
	formLayout2.AddRow5(button)

	horizontalLayout.AddLayout(formLayout, 0)
	horizontalLayout.AddLayout(formLayout2, 0)

	layout.AddLayout(horizontalLayout, 0)

	buttons := widgets.NewQDialogButtonBox(nil)
	buttons.SetOrientation(core.Qt__Horizontal)
	buttons.SetStandardButtons(widgets.QDialogButtonBox__Ok | widgets.QDialogButtonBox__Cancel)
	buttons.ConnectAccepted(func() {
		if passwordField.Text() == repeatField.Text() && usernameField.Text() != "" {
			dialog.Accept()
		} else {
			showError("Missing username or passwords do not match!")
		}
	})
	buttons.ConnectRejected(func() {
		dialog.Reject()
	})
	layout.AddWidget(buttons, 0, core.Qt__AlignRight)

	dialog.SetModal(true)
	dialog.Show()

	if dialog.Exec() == int(widgets.QDialog__Accepted) {
		return models.Secret{
			Title:       titleField.Text(),
			Username:    usernameField.Text(),
			Password:    passwordField.Text(),
			URL:         urlField.Text(),
			Description: descriptionField.ToPlainText(),
		}
	}
	return models.Secret{}
}

func saveFile() string {
	dialog := widgets.NewQFileDialog(nil, 0)
	file := dialog.GetSaveFileName(nil, "Create new database", "", "Database (*.db)", "", 0)
	return file
}

func loadFile() string {
	dialog := widgets.NewQFileDialog(nil, 0)
	file := dialog.GetOpenFileName(nil, "Open database", "", "Database (*.db)", "", 0)
	return file
}

func showError(message string) {
	dialog := widgets.NewQMessageBox(nil)
	dialog.SetWindowTitle("Error")
	dialog.SetText(message)
	dialog.SetIcon(widgets.QMessageBox__Critical)
	dialog.SetStandardButtons(widgets.QMessageBox__Ok)
	dialog.SetDefaultButton2(widgets.QMessageBox__Ok)
	dialog.SetEscapeButton2(widgets.QMessageBox__Ok)
	dialog.SetModal(true)
	dialog.Show()
	dialog.Exec()
}

func main() {
	log.Println("Start application")
	app := widgets.NewQApplication(len(os.Args), os.Args)
	window := widgets.NewQMainWindow(nil, 0)
	icon := gui.NewQIcon5("icons/pepega.png")
	window.SetWindowIcon(icon)
	window.SetMinimumSize2(800, 600)
	window.SetWindowTitle("Password manager")
	menu := createMenu()
	window.SetMenuBar(menu)
	tool := createToolBar()
	window.AddToolBar(core.Qt__TopToolBarArea, tool)
	central := widgets.NewQWidget(nil, 0)
	mainLayout := widgets.NewQVBoxLayout2(central)
	side := createSideMenu()
	main := createMain()
	splitter := widgets.NewQSplitter(nil)
	splitter.AddWidget(side)
	splitter.AddWidget(main)
	mainLayout.AddWidget(splitter, 0, 0)
	window.SetCentralWidget(central)
	window.Show()
	app.Exec()
}
