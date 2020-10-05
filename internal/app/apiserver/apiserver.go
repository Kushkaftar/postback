package apiserver

import (
	"fmt"
	"io"
	"io/ioutil"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
)

// APIServer ...
type APIServer struct {
	config *Config
	logger *logrus.Logger
	router *mux.Router
}

// New ...
func New(config *Config) *APIServer {
	return &APIServer{
		config: config,
		logger: logrus.New(),
		router: mux.NewRouter(),
	}
}

// Start ...
func (s *APIServer) Start() error {
	if err := s.configureLogger(); err != nil {
		return err
	}

	s.configureRouter()

	s.logger.Info("starting api server")
	return http.ListenAndServe(s.config.BindAddr, s.router)
}

func (s *APIServer) configureLogger() error {
	level, err := logrus.ParseLevel(s.config.LogLevel)
	if err != nil {
		return err
	}

	s.logger.SetLevel(level)

	return nil
}

func (s *APIServer) configureRouter() {
	s.router.HandleFunc("/hello", s.handleHello())
	s.router.HandleFunc("/postbackGet", s.postbackGet()).Methods("GET")
	s.router.HandleFunc("/postbackPost", s.postbackPost()).Methods("POST")
}

func (s *APIServer) handleHello() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		io.WriteString(w, "Hello")
	}
}

func (s *APIServer) postbackGet() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		bodys := r.URL.Scheme
		fmt.Println(bodys)
		requestURI := r.URL.Query().Encode()
		s.logger.Info(requestURI)

		io.WriteString(w, r.URL.RawQuery)
	}
}

func (s *APIServer) postbackPost() http.HandlerFunc {
	type Reseller struct {
		OfferID string `json:"offer_id"`
		UserID  string `json:"user_id"`
	}
	return func(w http.ResponseWriter, r *http.Request) {
		// g := Reseller{}
		// //fmt.Println(URLRESELLER)
		// err := json.NewDecoder(r.Body).Decode(&g)
		// s.logger.Error(err)

		// if err != nil {
		// 	w.WriteHeader(http.StatusBadRequest)
		// 	w.Write([]byte(`{"error":"` + err.Error() + `"`))
		// 	return
		// }

		body, err := ioutil.ReadAll(r.Body)
		fmt.Printf("%T\n", body)
		if err != nil {
			s.logger.Error(err)
		}
		// json.Unmarshal(body, &g)

		// fmt.Println(r)

		s.logger.Info(string(body))
		io.WriteString(w, r.PostForm.Encode())
	}
}
