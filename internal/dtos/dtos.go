package dtos

import "sica/internal/models"

type ProductSimpleResponse struct {
	ID          uint    `json:"id"`
	Name        string  `json:"name"`
	Price       float64 `json:"price"`
	Available	bool	`json:"available"`
	Description string  `json:"description"`
	Image       string  `json:"image"`
	Category    string  `json:"category"`
}

type ProductDetailsResponse struct {
	ID          uint    `json:"id"`
	Name        string  `json:"name"`
	Price       float64 `json:"price"`
	Description string  `json:"description"`
	Availabe    bool    `json:"available"`
	Visible     bool    `json:"visible"`
	Image       string  `json:"image"`
	Category    uint  `json:"category"`
	CategoryName string `json:"category_name"`
}

type CategoryDetailsResponse struct {
	Name     string                  `json:"name"`
	Order    uint                    `json:"order"`
	Products []ProductSimpleResponse `json:"products"`
}

type CategorySimpleResponse struct {
	ID    uint   `json:"id"`
	Name  string `json:"name"`
	Order uint   `json:"order"`
}

func ToCategoriesDetailsResponse(categories []models.Category) []CategoryDetailsResponse {

	//main page and product details

	var category []CategoryDetailsResponse

	for _, c := range categories {
		var product []ProductSimpleResponse
		for _, p := range c.Products {
			product = append(product, ProductSimpleResponse{
				ID:          p.ID,
				Name:        p.Name,
				Price:       p.Price,
				Description: p.Description,
				Available:   p.Available,
				Image:       p.Image,
				Category:    c.Name,
			})
		}

		category = append(category, CategoryDetailsResponse{
			Name:     c.Name,
			Order:    c.Order,
			Products: product,
		})
	}

	return category
}

func ToCategoriesSimpleResponse(categories []models.Category) []CategorySimpleResponse {

	//Manage categories

	var category []CategorySimpleResponse

	for _, c := range categories {
		category = append(category, CategorySimpleResponse{
			ID:    c.ID,
			Name:  c.Name,
			Order: c.Order,
		})
	}

	return category
}

func ToCategorySimpleResponse(category *models.Category) CategorySimpleResponse {

	//Create and Update categories

	return CategorySimpleResponse{
		ID:    category.ID,
		Name:  category.Name,
		Order: category.Order,
	}
}

func ToProductDetailsResponse(product *models.Product) ProductDetailsResponse{

	//Create and Update products

	return ProductDetailsResponse{
		ID: product.ID,
		Name: product.Name,
		Description : product.Description,
		Price: product.Price,
		Category: product.CategoryID,
		CategoryName: product.Category.Name,
		Availabe: product.Available,
		Visible: product.Visible,
		Image: product.Image,
	}
}

func ToProductsDetailsResponse(products []models.Product) []ProductDetailsResponse {

	//Manage products

	var dtoProducts []ProductDetailsResponse

	for _, p := range products {
		dtoProducts = append(dtoProducts, ProductDetailsResponse{
			ID: p.ID,
			Name: p.Name,
			Description : p.Description,
			Price: p.Price,
			Category: p.CategoryID,
			CategoryName: p.Category.Name,
			Availabe: p.Available,
			Visible: p.Visible,
			Image: p.Image,
		})
	}

	return dtoProducts
}





/*  {
  "message": "product created",
  "product": {
    "id": 4,
    "name": "Teclado mecánico RGB",
    "description": "Teclado mecánico con retroiluminación RGB y switches Cherry MX Red.",
    "price": 79.99,
    "category": 2,
    "category_detail": {
      "id": 2,
      "name": "Electronica",
      "order": 0,
      "products": null
    },
    "available": true,
    "visible": true,
    "image": "https://example.com/images/teclado-rgb.jpg"
  }
} */