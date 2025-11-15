package api

import (
	"fmt"
	"log"
	"net/http"
	"rsoi/internal/model"
	"rsoi/internal/repo"
	"rsoi/internal/service"
	"strconv"

	"github.com/gin-gonic/gin"
)

type RestApi struct {
	router  *gin.Engine
	service *service.Service
}

func NewApi(service *service.Service) *RestApi {
	return &RestApi{
		router:  gin.Default(),
		service: service,
	}
}

func (r *RestApi) Init() {
	v1 := r.router.Group("/api/v1")
	personsGroup := v1.Group("/persons")
	personsGroup.GET("", r.listPersons)
	personsGroup.POST("", r.createPerson)
	personsGroup.GET("/:id", r.getPerson)
	personsGroup.PATCH("/:id", r.updatePerson)
	personsGroup.DELETE("/:id", r.deletePerson)
}

func (r *RestApi) Run(addr string) error {
	return r.router.Run(addr)
}

func (r *RestApi) listPersons(c *gin.Context) {
	ps, err := r.service.GetAll()
	if err != nil {
		c.Status(http.StatusInternalServerError)
		return
	}
	if ps == nil {
		ps = []model.Person{}
	}
	c.JSON(http.StatusOK, ps)
}

func (r *RestApi) createPerson(c *gin.Context) {
	var dto CreatePersonDTO

	if err := c.ShouldBindJSON(&dto); err != nil {
		c.JSON(http.StatusBadRequest, err)
		return
	}

	person := model.Person{
		Name:    dto.Name,
		Age:     dto.Age,
		Address: dto.Address,
		Work:    dto.Work,
	}
	err := r.service.Add(&person)
	if err != nil {
		c.Status(http.StatusInternalServerError)
		return
	}

	c.Header("Location", fmt.Sprintf("/api/persons/%d", person.ID))
	c.Status(http.StatusCreated)
}

func (r *RestApi) getPerson(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil || id == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID format"})
		return
	}

	p, err := r.service.GetById(id)
	if err == repo.ErrNotFound {
		c.JSON(http.StatusNotFound, gin.H{"error": "Person not found"})
		return
	}

	c.JSON(http.StatusOK, p)
}

func (r *RestApi) updatePerson(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil || id == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	var dto UpdatePersonDTO
	if err := c.ShouldBindJSON(&dto); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	person := model.Person{
		ID:      id,
		Name:    dto.Name,
		Age:     dto.Age,
		Address: dto.Address,
		Work:    dto.Work,
	}

	err = r.service.Update(&person)
	if err != nil {
		if err == repo.ErrNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "Person not found"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}

	c.JSON(http.StatusOK, person)
}

func (r *RestApi) deletePerson(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil || id == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID format"})
		return
	}

	if err := r.service.Delete(id); err != nil {
		if err == repo.ErrNotFound {
			// c.JSON(http.StatusNotFound, gin.H{"error": "Person not found"})
			c.Status(http.StatusNoContent)
		} else {
			log.Println("ошибка при удалении:", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}

	c.Status(http.StatusNoContent)
}
