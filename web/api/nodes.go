/*
 Copyright 2018 Padduck, LLC
 Licensed under the Apache License, Version 2.0 (the "License");
 you may not use this file except in compliance with the License.
 You may obtain a copy of the License at
 	http://www.apache.org/licenses/LICENSE-2.0
 Unless required by applicable law or agreed to in writing, software
 distributed under the License is distributed on an "AS IS" BASIS,
 WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 See the License for the specific language governing permissions and
 limitations under the License.
*/

package api

import (
	"github.com/gin-gonic/gin"
	"github.com/pufferpanel/pufferpanel/models"
	"github.com/pufferpanel/pufferpanel/models/view"
	"github.com/pufferpanel/pufferpanel/services"
	"github.com/pufferpanel/pufferpanel/shared"
	builder "github.com/pufferpanel/apufferi/http"
	"net/http"
	"strconv"
)

func registerNodes(g *gin.RouterGroup) {
	g.Handle("GET", "", GetAllNodes)
	g.Handle("OPTIONS", "", shared.CreateOptions("GET"))

	g.Handle("PUT", "", CreateNode)
	g.Handle("GET", "/:id", GetNode)
	g.Handle("POST", "/:id", UpdateNode)
	g.Handle("DELETE", "/:id", DeleteNode)
	g.Handle("OPTIONS", "/:id", shared.CreateOptions("PUT", "GET", "POST", "DELETE"))
}

func GetAllNodes(c *gin.Context) {
	var ns *services.NodeService
	var err error
	response := builder.Respond(c)

	if ns, err = services.GetNodeService(); handleError(response, err) {
		return
	}

	var nodes *models.Nodes
	if nodes, err = ns.GetAll(); handleError(response, err) {
		return
	}

	data := view.FromNodes(nodes)

	response.Data(data).Send()
}

func GetNode(c *gin.Context) {
	var ns *services.NodeService
	var err error
	response := builder.Respond(c)

	if ns, err = services.GetNodeService(); handleError(response, err) {
		return
	}

	id, ok := validateId(c, response)
	if !ok {
		return
	}

	node, exists, err := ns.Get(id)
	if handleError(response, err) {
		return
	} else if !exists {
		response.Fail().Status(http.StatusNotFound).Message("no node with given id").Send()
		return
	}

	data := view.FromNode(node)

	response.Data(data).Send()
}

func CreateNode (c *gin.Context) {
	var ns *services.NodeService
	var err error
	response := builder.Respond(c)

	if ns, err = services.GetNodeService(); handleError(response, err) {
		return
	}

	model := view.NodeViewModel{}
	c.BindJSON(&model)

	create := &models.Node{}
	model.CopyToModel(create)
	if err = ns.Create(create); handleError(response, err) {
		return
	}

	response.Data(create).Send()
}

func UpdateNode (c *gin.Context) {
	var ns *services.NodeService
	var err error
	response := builder.Respond(c)

	if ns, err = services.GetNodeService(); handleError(response, err) {
		return
	}

	viewModel := &view.NodeViewModel{}
	c.BindJSON(viewModel)

	id, ok := validateId(c, response)
	if !ok {
		return
	}

	node, exists, err := ns.Get(id)
	if handleError(response, err) {
		return
	} else if !exists {
		response.Fail().Status(http.StatusNotFound).Message("no node with given id").Send()
		return
	}

	viewModel.CopyToModel(node)
	if err = ns.Update(node); handleError(response, err) {
		return
	}

	response.Data(node).Send()
}

func DeleteNode (c *gin.Context) {
	var ns *services.NodeService
	var err error
	response := builder.Respond(c)

	if ns, err = services.GetNodeService(); handleError(response, err) {
		return
	}

	id, ok := validateId(c, response)
	if !ok {
		return
	}

	node, exists, err := ns.Get(id)
	if handleError(response, err) {
		return
	} else if !exists {
		response.Fail().Status(http.StatusNotFound).Message("no node with given id").Send()
		return
	}

	err = ns.Delete(node.ID)
	if handleError(response, err) {
		return
	}

	response.Data(node).Send()
}

func validateId(c *gin.Context, response builder.Builder) (uint, bool) {
	param := c.Param("id")

	id, err := strconv.Atoi(param)

	if err != nil || id <= 0 {
		response.Fail().Status(http.StatusBadRequest).Message("id must be a positive number").Send()
		return 0, false
	}

	return uint(id), true
}

func handleError(response builder.Builder, err error) bool {
	if err != nil {
		response.Fail().Status(http.StatusInternalServerError).Message(err.Error()).Send()
		return true
	}
	return false
}