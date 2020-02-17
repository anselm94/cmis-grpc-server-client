# A Documents Server & Workbench Client

A 3 day implementation of a Documents Server and a Client Workbench allowing creating/deleting folders/documents in the server. Both the Documents server as well as Workbench are written in Golang. The metadata is persisted in PostgresqlDB. The server-client interface is through `gRPC` which is a *Remote-Procedure-Call* protocol.

## Overview

![Demo](./docs/demo.gif)

The *domain-modelling* follows that of *CMIS* but is minimal to show key features of Document Management systems like creating type definitions and property definitions, creating/deleting/fetching folders/documents with attached properties. The folders and documents can be attached to any custom type definitions and can have custom properties based on the property definitions, the type definitions have.

The server is written using Golang and utilizes the [concurrency](https://tour.golang.org/concurrency/1) for high-performance and supports server to client push via GORM's DB callbacks (which is not a DB trigger) after every create/delete of objects in DB.

The client is written in Golang too and has a peculiar Terminal UI (called *tui*). The client too has concurrency support while listening to server for updates, without blocking any threads. The UI runs in the main routine (similar to threads of Java world, but much simpler and robust, and sleeps while waiting for updates), while the server-listener runs in another go-routine.

![Architecture](./docs/architecture.png)

### Key Features:

* CMIS based domain modelling (minimal containing *Repository*, *Type Definition*, *Property Definition*, *Cmis Object*, *Cmis Property*, *Multifiling/Unfiling*) (See [dao.go](./internal/server/model/dao.go))

* Follows [3factor.app](https://3factor.app) design pattern (Create/Update/Delete actions are initiated by client -> server while Reads are pushed from server -> client)

* Solves the *Refresh* issue, where the objects are updated in the server and the client has to refreshed explicitly to get the updated information

* gRPC based client-server communication support full bidirectional streaming

* Just *3 days* to build from concept to implementation!

### Workbench Client

Workbench client is inspired by CMIS Workbench. This is a terminal based UI (commonly termed as *tui*) and is built using [tui-go](https://github.com/marcusolsson/tui-go). You can navigate the UI using ⬆️/⬇️ *arrow* keys and jump between the sections using ➡️ *tab* key and enter into a folder or trigger an action using ↩️ *enter/return* key.

![Client Workbench running in tui](./docs/client-workbench.png)

## How to Setup

Setup involves setting up `1` server and `n` clients where `n >= 0`. The server and clients interact via gRPC protocol. 

### Prerequisite

1) Clone the Git repository and navigate into the folder

```
git clone https://github.wdf.sap.corp/I327891/documents-server-workbench.git
cd documents-server-workbench
```

2) Install [Golang tools](https://golang.org/dl/) for your platform and update `PATH`, `GOPATH` environmental variables as necessary

3) [Install and setup PostgreSQL](https://www.postgresql.org/download/) locally or connect to any running local/cloud instance

3) Edit [Config.go](./Config.go) and update the values accordingly

### Setting up Documents Server

1) Start the gRPC based Documents Server

* *First Run* - Inorder to populate some sample data, run the server with `-populate` flag

    ```
    go run cmd/server/main.go -populate
    ```

* *Subsequent Runs* - If wanted to reuse the existing data in DB, run the server normally

    ```
    go run cmd/server/main.go
    ```

2) Press `ctrl + c` to terminate the server

### Setting up Client workbench

1) Open a new terminal session and start any number of gRPC based Client Workbench

```
go run cmd/client/main.go
```

2) Press `Esc` to terminate the UI

### How to Use

![How to use Client Workbench](./docs/client-workbench-instructions.png)

* Use the ➡️ *tab* keys to jump between the sections in the following order

  'Documents' area -> Folder Name -> Create Folder -> Document Name -> Create Document -> Delete Object

* Whie in *'Documents'* area, use ⬆️/⬇️ *arrow* keys to select a folder/document

* While in *'Documents'* area, press ↩️ *enter/return* key to navigate into/out of the folder

* While in *'Actions'* area, pressing ↩️ *enter/return* key after selecting any of the actions, triggers the action