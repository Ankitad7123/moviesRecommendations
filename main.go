package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

type Movie struct {
	Title      string `json:"title"`
	IMDbLink   string `json:"link,omitempty"`
	PosterURL  string `json:"poster_url,omitempty"`
	PosterPath string `json:"poster_path,omitempty"`
}

type TMDbAPIResponse struct {
	Results []Movie `json:"results"`
}

// TMDb API key (replace with your credentials)
const apiKey = ""

func main() {
	r := gin.Default()

	// Define routes
	r.GET("/recommendations", getMovieRecommendations)
	r.GET("/recommendations2", getMovieRecommendations2)
	r.GET("/indian-movies", getIndianMovies)
	r.GET("/upcoming-movies", getUpcomingMovies)
	r.GET("/most-anticipated-movies", getMostAnticipatedMovies)
	r.GET("/hollywood-movies", getHollywoodMovies)
	r.GET("/top-rated-movies", getTopRatedMovies)
	r.GET("/popular-movies", getPopularMovies)
	r.GET("/now-playing-movies", getNowPlayingMovies)
	r.GET("/movie-genres", getMovieGenres)

	// Start the server
	r.Run(":8080")
}

// Handler function to get movie recommendations from TMDb
func getMovieRecommendations(c *gin.Context) {
	movieTitle := "indian" // Change this to accept user input
	apiURL := fmt.Sprintf("https://api.themoviedb.org/3/search/movie?api_key=%s&query=%s", apiKey, movieTitle)

	resp, err := makeTMDbRequest(apiURL)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	var posterURLs []string
	for _, movie := range resp.Results {
		if movie.PosterPath != "" {
			posterURL := fmt.Sprintf("https://image.tmdb.org/t/p/w500%s", movie.PosterPath)
			posterURLs = append(posterURLs, posterURL)
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"poster_urls": posterURLs,
	})
}

func getMovieRecommendations2(c *gin.Context) {
	movieTitle := c.Query("title")
	if movieTitle == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Movie title is required"})
		return
	}

	movieID, err := getMovieID(movieTitle)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	recommendations, err := fetchRecommendations(movieID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"recommendations": recommendations,
	})
}

func getIndianMovies(c *gin.Context) {
	apiURL := "https://api.themoviedb.org/3/discover/movie?api_key=" + apiKey + "&with_original_language=hi"
	resp, err := makeTMDbRequest(apiURL)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, resp)
}

func getUpcomingMovies(c *gin.Context) {
	apiURL := "https://api.themoviedb.org/3/movie/upcoming?api_key=" + apiKey + "&language=en-US&page=1"
	resp, err := makeTMDbRequest(apiURL)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, resp)
}

func getMostAnticipatedMovies(c *gin.Context) {
	apiURL := "https://api.themoviedb.org/3/movie/popular?api_key=" + apiKey
	resp, err := makeTMDbRequest(apiURL)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, resp)
}

func getHollywoodMovies(c *gin.Context) {
	apiURL := "https://api.themoviedb.org/3/discover/movie?api_key=" + apiKey + "&with_original_language=en"
	resp, err := makeTMDbRequest(apiURL)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, resp)
}

func getTopRatedMovies(c *gin.Context) {
	apiURL := "https://api.themoviedb.org/3/movie/top_rated?api_key=" + apiKey + "&language=en-US&page=1"
	resp, err := makeTMDbRequest(apiURL)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, resp)
}

func getPopularMovies(c *gin.Context) {
	apiURL := "https://api.themoviedb.org/3/movie/popular?api_key=" + apiKey + "&language=en-US&page=1"
	resp, err := makeTMDbRequest(apiURL)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, resp)
}

func getNowPlayingMovies(c *gin.Context) {
	apiURL := "https://api.themoviedb.org/3/movie/now_playing?api_key=" + apiKey + "&language=en-US&page=1"
	resp, err := makeTMDbRequest(apiURL)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, resp)
}

func getMovieGenres(c *gin.Context) {
	apiURL := "https://api.themoviedb.org/3/genre/movie/list?api_key=" + apiKey + "&language=en-US"
	resp, err := makeTMDbRequest(apiURL)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, resp)
}

func makeTMDbRequest(url string) (TMDbAPIResponse, error) {
	resp, err := http.Get(url)
	if err != nil {
		return TMDbAPIResponse{}, fmt.Errorf("error making request to TMDb: %v", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return TMDbAPIResponse{}, fmt.Errorf("error reading response: %v", err)
	}

	var apiResponse TMDbAPIResponse
	if err := json.Unmarshal(body, &apiResponse); err != nil {
		return TMDbAPIResponse{}, fmt.Errorf("error parsing JSON: %v", err)
	}

	return apiResponse, nil
}

func getMovieID(title string) (int, error) {
	apiURL := fmt.Sprintf("https://api.themoviedb.org/3/search/movie?api_key=%s&query=%s", apiKey, strings.ReplaceAll(title, " ", "+"))

	resp, err := http.Get(apiURL)
	if err != nil {
		return 0, fmt.Errorf("error making request to TMDb: %v", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return 0, fmt.Errorf("error reading response: %v", err)
	}

	var searchResponse struct {
		Results []struct {
			ID int `json:"id"`
		} `json:"results"`
	}

	if err := json.Unmarshal(body, &searchResponse); err != nil {
		return 0, fmt.Errorf("error parsing JSON: %v", err)
	}

	if len(searchResponse.Results) == 0 {
		return 0, fmt.Errorf("no movie found with title: %s", title)
	}

	return searchResponse.Results[0].ID, nil
}

func fetchRecommendations(movieID int) ([]Movie, error) {
	apiURL := fmt.Sprintf("https://api.themoviedb.org/3/movie/%d/recommendations?api_key=%s", movieID, apiKey)

	resp, err := http.Get(apiURL)
	if err != nil {
		return nil, fmt.Errorf("error making request to TMDb: %v", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading response: %v", err)
	}

	var recommendationsResp struct {
		Results []Movie `json:"results"`
	}

	if err := json.Unmarshal(body, &recommendationsResp); err != nil {
		return nil, fmt.Errorf("error parsing JSON: %v", err)
	}

	return recommendationsResp.Results, nil
}
