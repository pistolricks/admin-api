package main

import (
	"github.com/pistolricks/admin-api/internal/data"
	"github.com/pistolricks/admin-api/internal/validator"
	"net/http"
)

// Using the data.Product type instead of defining our own

func (app *application) productsHandler(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Style        string
		ProductTitle string
		Mill         string
		CategoryName string
		ColorName    string
		data.Filters
	}

	v := validator.New()

	qs := r.URL.Query()

	input.Style = app.readString(qs, "attrs->>style", "")

	input.Mill = app.readString(qs, "attrs->>mill", "")

	input.Filters.Page = app.readInt(qs, "page", 1, v)
	input.Filters.PageSize = app.readInt(qs, "page_size", 20, v)
	input.Filters.Sort = app.readString(qs, "sort", "id")
	input.Filters.SortSafelist = []string{"id", "product_title", "category_name", "subcategory_name", "color_name", "sizes", "mill", "msrp", "suggested_price", "map_pricing", "-id", "-product_title", "-category_name", "-msrp", "-suggested_price", "-map_pricing", "-subcategory_name", "-color_name", "-size", "-mill"}

	if data.ValidateFilters(v, input.Filters); !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	products, metadata, err := app.models.Products.GetAll(input.Style, input.Mill, input.Filters)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	err = app.writeJSON(w, http.StatusOK, envelope{"products": products, "metadata": metadata}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *application) importProducts(w http.ResponseWriter, r *http.Request) {

	var input struct {
		Products []data.Attrs `json:"products"`
	}

	err := app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	for _, res := range input.Products {

		attrs := new(data.Attrs)

		attrs.ID = res.ID
		attrs.ProductTitle = res.ProductTitle
		attrs.ProductDescription = res.ProductDescription
		attrs.Style = res.Style
		attrs.AvailableSizes = res.AvailableSizes
		attrs.BrandLogoImage = res.BrandLogoImage
		attrs.ThumbnailImage = res.ThumbnailImage
		attrs.ColorSwatchImage = res.ColorSwatchImage
		attrs.ProductImage = res.ProductImage
		attrs.SpecSheet = res.SpecSheet
		attrs.PriceText = res.PriceText
		attrs.SuggestedPrice = res.SuggestedPrice
		attrs.CategoryName = res.CategoryName
		attrs.SubcategoryName = res.SubcategoryName
		attrs.ColorName = res.ColorName
		attrs.ColorSquareImage = res.ColorSquareImage
		attrs.ColorProductImage = res.ColorProductImage
		attrs.ColorProductImageThumbnail = res.ColorProductImageThumbnail
		attrs.Size = res.Size
		attrs.PieceWeight = res.PieceWeight
		attrs.PiecePrice = res.PiecePrice
		attrs.DozensPrice = res.DozensPrice
		attrs.CasePrice = res.CasePrice
		attrs.PriceGroup = res.PriceGroup
		attrs.CaseSize = res.CaseSize
		attrs.InventoryKey = res.InventoryKey
		attrs.SizeIndex = res.SizeIndex
		attrs.SanmarMainframeColor = res.SanmarMainframeColor
		attrs.Mill = res.Mill
		attrs.ProductStatus = res.ProductStatus
		attrs.CompanionStyle = res.CompanionStyle
		attrs.Msrp = res.Msrp
		attrs.MapPricing = res.MapPricing
		attrs.FrontModelImageUrl = res.FrontModelImageUrl
		attrs.BackModelImageUrl = res.BackModelImageUrl
		attrs.FrontFlatImageUrl = res.FrontFlatImageUrl
		attrs.BackFlatImageUrl = res.BackFlatImageUrl
		attrs.ProductMeasurements = res.ProductMeasurements
		attrs.PmsColor = res.PmsColor
		attrs.Gtin = res.Gtin + "GTIN"
		attrs.DecorationSpecSheet = res.DecorationSpecSheet

		product := new(data.Product)

		product.ID = res.ID
		product.Type = "sanmar"
		product.Attrs = *attrs

		err = app.models.Products.Insert(product)
		if err != nil {
			app.serverErrorResponse(w, r, err)
			return
		}

	}

}
