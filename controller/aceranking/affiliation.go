package aceranking

import (
	"acemap/controller"
	"acemap/data"
	"acemap/model/aceranking"
	"acemap/model/dao"
	"github.com/gin-gonic/gin"
	"net/http"
	"sort"
)

type affiliationForm struct {
	Type          string  `form:"type"`
	VenueIDs      string  `form:"venue_ids"`
	StartYear     int     `form:"start_year"`
	EndYear       int     `form:"end_year"`
	OrderBy       int     `form:"order_by"`
	FirstAuthor   int     `form:"first_author"`
	AffiliationID data.ID `form:"affiliation_id"`
}

func Affiliation(c *gin.Context) {
	var param affiliationForm
	var err error
	err = c.Bind(&param)
	if err != nil {
		c.JSON(http.StatusOK, controller.RespondError(err.Error()))
		return
	}

	venueMask, err := aceranking.GetVenueMask(param.VenueIDs)
	if err != nil {
		c.JSON(http.StatusOK, controller.RespondError(err.Error()))
		return
	}
	affIndex := data.AffiliationIDToIndex[param.AffiliationID]
	var details aceranking.AffiliationDetails
	if param.FirstAuthor == 0 {
		details = aceranking.GetAffDetailsOnAllAuthors(param.Type, venueMask, param.StartYear, param.EndYear, affIndex)
	} else if param.FirstAuthor == 1 {
		details = aceranking.GetAffDetailsOnFirstAuthorWeak(param.Type, venueMask, param.StartYear, param.EndYear, affIndex)
	} else {
		details = aceranking.GetAffDetailsOnFirstAuthorStrong(param.Type, venueMask, param.StartYear, param.EndYear, affIndex)
	}

	statistics := details.AuthorList
	aceranking.Sort(statistics, param.OrderBy)

	//end := 50
	//if len(statistics) < end {
	//	end = len(statistics)
	//}

	// Number of rows to show AT MOST.
	end := 50
	if len(statistics) < end {
		end = len(statistics)
	}
	// Order by. 1: count; 2: citation, 3: citation_share, 4: H-index, 5: AceScore
	switch param.OrderBy {
	case 1:
		sort.Slice(statistics, func(i, j int) bool {
			return statistics[i].Count > statistics[j].Count
		})
		for end > 0 && statistics[end-1].Count == 0 {
			end--
		}
	case 2:
		sort.Slice(statistics, func(i, j int) bool {
			return statistics[i].Citation > statistics[j].Citation
		})
		for end > 0 && statistics[end-1].Citation == 0 {
			end--
		}
	case 3:
		sort.Slice(statistics, func(i, j int) bool {
			return statistics[i].CitationShare > statistics[j].CitationShare
		})
		for end > 0 && statistics[end-1].CitationShare == 0 {
			end--
		}
	case 4:
		sort.Slice(statistics, func(i, j int) bool {
			return statistics[i].HIndex > statistics[j].HIndex
		})
		for end > 0 && statistics[end-1].HIndex == 0 {
			end--
		}
	case 5:
		sort.Slice(statistics, func(i, j int) bool {
			return statistics[i].AceScore > statistics[j].AceScore
		})
		for end > 0 && statistics[end-1].AceScore == 0 {
			end--
		}
	}

	for i := 0; i < end; i++ {
		statistics[i].AuthorName = dao.GetAuthorNameByID(statistics[i].AuthorID)
	}
	details.AuthorList = statistics[:end]

	c.JSON(http.StatusOK, controller.RespondData(details))
}
