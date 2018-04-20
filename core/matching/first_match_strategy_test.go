package matching_test

import (
	"testing"

	"github.com/SpectoLabs/hoverfly/core/matching"
	"github.com/SpectoLabs/hoverfly/core/models"
	. "github.com/SpectoLabs/hoverfly/core/util"
	. "github.com/onsi/gomega"
)

var testResponse = models.ResponseDetails{
	Body: "request matched",
}

func Test_FirstMatchStrategy_EmptyRequestMatchersShouldMatchOnAnyRequest(t *testing.T) {
	RegisterTestingT(t)

	simulation := models.NewSimulation()

	simulation.AddRequestMatcherResponsePair(&models.RequestMatcherResponsePair{
		RequestMatcher: models.RequestMatcher{},
		Response:       testResponse,
	})

	r := models.RequestDetails{
		Method:      "GET",
		Destination: "somehost.com",
		Headers: map[string][]string{
			"sdv": {"ascd"},
		},
	}
	result := matching.FirstMatchStrategy(r, false, simulation, make(map[string]string))

	Expect(result.Pair.Response.Body).To(Equal("request matched"))
}

func Test_FirstMatchStrategy_RequestMatchersShouldMatchOnBody(t *testing.T) {
	RegisterTestingT(t)

	simulation := models.NewSimulation()

	simulation.AddRequestMatcherResponsePair(&models.RequestMatcherResponsePair{
		RequestMatcher: models.RequestMatcher{
			Body: &models.RequestFieldMatchers{
				ExactMatch: StringToPointer("body"),
			},
		},
		Response: testResponse,
	})

	r := models.RequestDetails{
		Body: "body",
	}
	result := matching.FirstMatchStrategy(r, false, simulation, make(map[string]string))
	Expect(result.Error).To(BeNil())

	Expect(result.Pair.Response.Body).To(Equal("request matched"))
}

func Test_FirstMatchStrategy_ReturnResponseWhenAllHeadersMatch(t *testing.T) {
	RegisterTestingT(t)

	headers := map[string][]string{
		"header1": []string{"val1"},
		"header2": []string{"val2"},
	}

	simulation := models.NewSimulation()

	simulation.AddRequestMatcherResponsePair(&models.RequestMatcherResponsePair{
		RequestMatcher: models.RequestMatcher{
			Headers: headers,
		},
		Response: testResponse,
	})

	r := models.RequestDetails{
		Method:      "GET",
		Destination: "http://somehost.com",
		Headers: map[string][]string{
			"header1": []string{"val1"},
			"header2": []string{"val2"},
		},
	}

	result := matching.FirstMatchStrategy(r, false, simulation, make(map[string]string))

	Expect(result.Pair.Response.Body).To(Equal("request matched"))
}

func Test_FirstMatchStrategy_ReturnNilWhenOneHeaderNotPresentInRequest(t *testing.T) {
	RegisterTestingT(t)

	headers := map[string][]string{
		"header1": []string{"val1"},
		"header2": []string{"val2"},
	}

	simulation := models.NewSimulation()

	simulation.AddRequestMatcherResponsePair(&models.RequestMatcherResponsePair{
		RequestMatcher: models.RequestMatcher{
			Headers: headers,
		},
		Response: testResponse,
	})

	r := models.RequestDetails{
		Method:      "GET",
		Destination: "http://somehost.com",
		Headers: map[string][]string{
			"header1": []string{"val1"},
		},
	}

	result := matching.FirstMatchStrategy(r, false, simulation, make(map[string]string))

	Expect(result.Pair).To(BeNil())
}

func Test_FirstMatchStrategy_ReturnNilWhenOneHeaderValueDifferent(t *testing.T) {
	RegisterTestingT(t)

	headers := map[string][]string{
		"header1": []string{"val1"},
		"header2": []string{"val2"},
	}

	simulation := models.NewSimulation()

	simulation.AddRequestMatcherResponsePair(&models.RequestMatcherResponsePair{
		RequestMatcher: models.RequestMatcher{
			Headers: headers,
		},
		Response: testResponse,
	})

	r := models.RequestDetails{
		Method:      "GET",
		Destination: "somehost.com",
		Headers: map[string][]string{
			"header1": []string{"val1"},
			"header2": []string{"different"},
		},
	}
	result := matching.FirstMatchStrategy(r, false, simulation, make(map[string]string))

	Expect(result.Pair).To(BeNil())
}

func Test_FirstMatchStrategy_ReturnResponseWithMultiValuedHeaderMatch(t *testing.T) {
	RegisterTestingT(t)

	headers := map[string][]string{
		"header1": []string{"val1-a", "val1-b"},
		"header2": []string{"val2"},
	}

	simulation := models.NewSimulation()

	simulation.AddRequestMatcherResponsePair(&models.RequestMatcherResponsePair{
		RequestMatcher: models.RequestMatcher{
			Headers: headers,
		},
		Response: testResponse,
	})

	r := models.RequestDetails{
		Method:      "GET",
		Destination: "http://somehost.com",
		Body:        "test-body",
		Headers: map[string][]string{
			"header1": []string{"val1-a", "val1-b"},
			"header2": []string{"val2"},
		},
	}
	result := matching.FirstMatchStrategy(r, false, simulation, make(map[string]string))

	Expect(result.Pair.Response.Body).To(Equal("request matched"))
}

func Test_FirstMatchStrategy_ReturnNilWithDifferentMultiValuedHeaders(t *testing.T) {
	RegisterTestingT(t)

	headers := map[string][]string{
		"header1": []string{"val1-a", "val1-b"},
		"header2": []string{"val2"},
	}

	simulation := models.NewSimulation()

	simulation.AddRequestMatcherResponsePair(&models.RequestMatcherResponsePair{
		RequestMatcher: models.RequestMatcher{
			Headers: headers,
		},
		Response: testResponse,
	})

	r := models.RequestDetails{
		Method:      "GET",
		Destination: "http://somehost.com",
		Headers: map[string][]string{
			"header1": []string{"val1-a", "val1-differnet"},
			"header2": []string{"val2"},
		},
	}

	result := matching.FirstMatchStrategy(r, false, simulation, make(map[string]string))

	Expect(result.Pair).To(BeNil())
}

func Test_FirstMatchStrategy_EndpointMatchWithHeaders(t *testing.T) {
	RegisterTestingT(t)

	headers := map[string][]string{
		"header1": []string{"val1-a", "val1-b"},
		"header2": []string{"val2"},
	}

	destination := "testhost.com"
	method := "GET"
	path := "/a/1"
	query := "q=test"

	simulation := models.NewSimulation()

	simulation.AddRequestMatcherResponsePair(&models.RequestMatcherResponsePair{
		RequestMatcher: models.RequestMatcher{
			Headers: headers,
			Destination: &models.RequestFieldMatchers{
				Matcher:    "exact",
				Value:      "testhost.com",
				ExactMatch: &destination,
			},
			Path: &models.RequestFieldMatchers{
				Matcher:    "exact",
				Value:      "/a/1",
				ExactMatch: &path,
			},
			Method: &models.RequestFieldMatchers{
				Matcher:    "exact",
				Value:      "GET",
				ExactMatch: &method,
			},
			Query: &models.RequestFieldMatchers{
				Matcher:    "exact",
				Value:      "q=test",
				ExactMatch: &query,
			},
		},
		Response: testResponse,
	})

	r := models.RequestDetails{
		Method:      "GET",
		Destination: "testhost.com",
		Path:        "/a/1",
		Query: map[string][]string{
			"q": []string{"test"},
		},
		Headers: map[string][]string{
			"header1": []string{"val1-a", "val1-b"},
			"header2": []string{"val2"},
		},
	}
	result := matching.FirstMatchStrategy(r, false, simulation, make(map[string]string))

	Expect(result.Pair.Response.Body).To(Equal("request matched"))
}

func Test_FirstMatchStrategy_EndpointMismatchWithHeadersReturnsNil(t *testing.T) {
	RegisterTestingT(t)

	headers := map[string][]string{
		"header1": []string{"val1-a", "val1-b"},
		"header2": []string{"val2"},
	}

	destination := "testhost.com"
	method := "GET"
	path := "/a/1"
	query := "q=test"

	simulation := models.NewSimulation()

	simulation.AddRequestMatcherResponsePair(&models.RequestMatcherResponsePair{
		RequestMatcher: models.RequestMatcher{
			Headers: headers,
			Destination: &models.RequestFieldMatchers{
				Matcher:    "exact",
				Value:      "testhost.com",
				ExactMatch: &destination,
			},
			Path: &models.RequestFieldMatchers{
				ExactMatch: &path,
			},
			Method: &models.RequestFieldMatchers{
				ExactMatch: &method,
			},
			Query: &models.RequestFieldMatchers{
				ExactMatch: &query,
			},
		},
		Response: testResponse,
	})

	r := models.RequestDetails{
		Method:      "GET",
		Destination: "http://testhost.com",
		Path:        "/a/1",
		Query: map[string][]string{
			"q": []string{"different"},
		},
		Headers: map[string][]string{
			"header1": []string{"val1-a", "val1-b"},
			"header2": []string{"val2"},
		},
	}

	result := matching.FirstMatchStrategy(r, false, simulation, make(map[string]string))

	Expect(result.Pair).To(BeNil())
}

func Test_FirstMatchStrategy_AbleToMatchAnEmptyPathInAReasonableWay(t *testing.T) {
	RegisterTestingT(t)

	destination := "testhost.com"
	method := "GET"
	path := ""
	query := "q=test"
	simulation := models.NewSimulation()

	simulation.AddRequestMatcherResponsePair(&models.RequestMatcherResponsePair{
		RequestMatcher: models.RequestMatcher{
			Destination: &models.RequestFieldMatchers{
				Matcher:    "exact",
				Value:      "testhost.com",
				ExactMatch: &destination,
			},
			Path: &models.RequestFieldMatchers{
				Matcher:    "exact",
				Value:      "",
				ExactMatch: &path,
			},
			Method: &models.RequestFieldMatchers{
				Matcher:    "exact",
				Value:      "GET",
				ExactMatch: &method,
			},
			Query: &models.RequestFieldMatchers{
				Matcher:    "exact",
				Value:      "q=test",
				ExactMatch: &query,
			},
		},
		Response: testResponse,
	})

	r := models.RequestDetails{
		Method:      "GET",
		Destination: "testhost.com",
		Query: map[string][]string{
			"q": []string{"test"},
		},
	}
	result := matching.FirstMatchStrategy(r, false, simulation, make(map[string]string))

	Expect(result.Pair.Response.Body).To(Equal("request matched"))

	r = models.RequestDetails{
		Method:      "GET",
		Destination: "testhost.com",
		Path:        "/a/1",
		Query: map[string][]string{
			"q": []string{"test"},
		},
	}

	result = matching.FirstMatchStrategy(r, false, simulation, make(map[string]string))

	Expect(result.Pair).To(BeNil())
}

func Test_FirstMatchStrategy_RequestMatcherResponsePairCanBeConvertedToARequestResponsePairView_WhileIncomplete(t *testing.T) {
	RegisterTestingT(t)

	method := "POST"

	requestMatcherResponsePair := models.RequestMatcherResponsePair{
		RequestMatcher: models.RequestMatcher{
			Method: &models.RequestFieldMatchers{
				ExactMatch: &method,
			},
		},
		Response: testResponse,
	}

	pairView := requestMatcherResponsePair.BuildView()

	Expect(pairView.RequestMatcher.Method.ExactMatch).To(Equal(StringToPointer("POST")))
	Expect(pairView.RequestMatcher.Destination).To(BeNil())
	Expect(pairView.RequestMatcher.Path).To(BeNil())
	Expect(pairView.RequestMatcher.Scheme).To(BeNil())
	Expect(pairView.RequestMatcher.Query).To(BeNil())

	Expect(pairView.Response.Body).To(Equal("request matched"))
}

func Test_FirstMatchStrategy_RequestMatchersCanUseGlobsAndBeMatched(t *testing.T) {
	RegisterTestingT(t)

	simulation := models.NewSimulation()

	simulation.AddRequestMatcherResponsePair(&models.RequestMatcherResponsePair{
		RequestMatcher: models.RequestMatcher{
			Destination: &models.RequestFieldMatchers{
				GlobMatch: StringToPointer("*.com"),
			},
		},
		Response: testResponse,
	})

	request := models.RequestDetails{
		Method:      "GET",
		Destination: "testhost.com",
		Path:        "/api/1",
	}

	result := matching.FirstMatchStrategy(request, false, simulation, make(map[string]string))
	Expect(result.Error).To(BeNil())

	Expect(result.Pair.Response.Body).To(Equal("request matched"))
}

func Test_FirstMatchStrategy_RequestMatchersCanUseGlobsOnSchemeAndBeMatched(t *testing.T) {
	RegisterTestingT(t)

	simulation := models.NewSimulation()

	simulation.AddRequestMatcherResponsePair(&models.RequestMatcherResponsePair{
		RequestMatcher: models.RequestMatcher{
			Scheme: &models.RequestFieldMatchers{
				GlobMatch: StringToPointer("H*"),
			},
		},
		Response: testResponse,
	})

	request := models.RequestDetails{
		Method:      "GET",
		Destination: "testhost.com",
		Scheme:      "http",
		Path:        "/api/1",
	}

	result := matching.FirstMatchStrategy(request, false, simulation, make(map[string]string))
	Expect(result.Error).To(BeNil())

	Expect(result.Pair.Response.Body).To(Equal("request matched"))
}

func Test_FirstMatchStrategy_RequestMatchersCanUseGlobsOnHeadersAndBeMatched(t *testing.T) {
	RegisterTestingT(t)

	simulation := models.NewSimulation()

	simulation.AddRequestMatcherResponsePair(&models.RequestMatcherResponsePair{
		RequestMatcher: models.RequestMatcher{
			Headers: map[string][]string{
				"unique-header": []string{"*"},
			},
		},
		Response: testResponse,
	})

	request := models.RequestDetails{
		Method:      "GET",
		Destination: "testhost.com",
		Path:        "/api/1",
		Headers: map[string][]string{
			"unique-header": []string{"totally-unique"},
		},
	}

	result := matching.FirstMatchStrategy(request, false, simulation, make(map[string]string))
	Expect(result.Error).To(BeNil())

	Expect(result.Pair.Response.Body).To(Equal("request matched"))
}

func Test_FirstMatchStrategy_RequestMatcherResponsePair_ConvertToRequestResponsePairView_CanBeConvertedToARequestResponsePairView_WhileIncomplete(t *testing.T) {
	RegisterTestingT(t)

	method := "POST"

	requestMatcherResponsePair := models.RequestMatcherResponsePair{
		RequestMatcher: models.RequestMatcher{
			Method: &models.RequestFieldMatchers{
				ExactMatch: &method,
			},
		},
		Response: testResponse,
	}

	pairView := requestMatcherResponsePair.BuildView()

	Expect(pairView.RequestMatcher.Method.ExactMatch).To(Equal(StringToPointer("POST")))
	Expect(pairView.RequestMatcher.Destination).To(BeNil())
	Expect(pairView.RequestMatcher.Path).To(BeNil())
	Expect(pairView.RequestMatcher.Scheme).To(BeNil())
	Expect(pairView.RequestMatcher.Query).To(BeNil())

	Expect(pairView.Response.Body).To(Equal("request matched"))
}

func Test_FirstMatchShouldNotBeCachableIfMatchedOnEverythingApartFromHeadersAtLeastOnce(t *testing.T) {
	RegisterTestingT(t)

	simulation := models.NewSimulation()

	simulation.AddRequestMatcherResponsePair(&models.RequestMatcherResponsePair{
		RequestMatcher: models.RequestMatcher{
			Method: &models.RequestFieldMatchers{
				Matcher:    "exact",
				Value:      "POST",
				ExactMatch: StringToPointer("POST"),
			},
			Body: &models.RequestFieldMatchers{
				Matcher:    "exact",
				Value:      "body",
				ExactMatch: StringToPointer("body"),
			},
			Scheme: &models.RequestFieldMatchers{
				Matcher:    "exact",
				Value:      "http",
				ExactMatch: StringToPointer("http"),
			},
			Query: &models.RequestFieldMatchers{
				Matcher:    "exact",
				Value:      "foo=bar",
				ExactMatch: StringToPointer("foo=bar"),
			},
			Path: &models.RequestFieldMatchers{
				Matcher:    "exact",
				Value:      "/foo",
				ExactMatch: StringToPointer("/foo"),
			},
			Destination: &models.RequestFieldMatchers{
				Matcher:    "exact",
				Value:      "www.test.com",
				ExactMatch: StringToPointer("www.test.com"),
			},
			Headers: map[string][]string{
				"foo": {"bar"},
			},
		},
		Response: testResponse,
	})

	simulation.AddRequestMatcherResponsePair(&models.RequestMatcherResponsePair{
		RequestMatcher: models.RequestMatcher{
			Method: &models.RequestFieldMatchers{
				Matcher:    "exact",
				Value:      "GET",
				ExactMatch: StringToPointer("GET"),
			},
		},
		Response: testResponse,
	})

	r := models.RequestDetails{
		Method:      "POST",
		Destination: "www.test.com",
		Query: map[string][]string{
			"foo": []string{"bar"},
		},
		Scheme: "http",
		Body:   "body",
		Path:   "/foo",
		Headers: map[string][]string{
			"miss": {"me"},
		},
	}

	result := matching.FirstMatchStrategy(r, false, simulation, make(map[string]string))

	Expect(result.Error).ToNot(BeNil())
	Expect(result.Cachable).To(BeFalse())
}

func Test_FirstMatchShouldBeCachableIfMatchedOnEverythingApartFromHeadersZeroTimes(t *testing.T) {
	RegisterTestingT(t)

	simulation := models.NewSimulation()

	simulation.AddRequestMatcherResponsePair(&models.RequestMatcherResponsePair{
		RequestMatcher: models.RequestMatcher{
			Method: &models.RequestFieldMatchers{
				Matcher:    "exact",
				Value:      "POST",
				ExactMatch: StringToPointer("POST"),
			},
			Body: &models.RequestFieldMatchers{
				Matcher:    "exact",
				Value:      "body",
				ExactMatch: StringToPointer("body"),
			},
			Scheme: &models.RequestFieldMatchers{
				Matcher:    "exact",
				Value:      "http",
				ExactMatch: StringToPointer("http"),
			},
			Query: &models.RequestFieldMatchers{
				Matcher:    "exact",
				Value:      "?foo=bar",
				ExactMatch: StringToPointer("?foo=bar"),
			},
			Path: &models.RequestFieldMatchers{
				Matcher:    "exact",
				Value:      "/foo",
				ExactMatch: StringToPointer("/foo"),
			},
			Destination: &models.RequestFieldMatchers{
				Matcher:    "exact",
				Value:      "www.test.com",
				ExactMatch: StringToPointer("www.test.com"),
			},
			Headers: map[string][]string{
				"foo": {"bar"},
			},
		},
		Response: testResponse,
	})

	simulation.AddRequestMatcherResponsePair(&models.RequestMatcherResponsePair{
		RequestMatcher: models.RequestMatcher{
			Method: &models.RequestFieldMatchers{
				Matcher:    "exact",
				Value:      "GET",
				ExactMatch: StringToPointer("GET"),
			},
		},
		Response: testResponse,
	})

	r := models.RequestDetails{
		Method:      "MISS",
		Destination: "www.test.com",
		Query: map[string][]string{
			"foo": []string{"bar"},
		},
		Scheme: "http",
		Body:   "body",
		Path:   "/foo",
		Headers: map[string][]string{
			"miss": {"me"},
		},
	}

	result := matching.FirstMatchStrategy(r, false, simulation, make(map[string]string))

	Expect(result.Error).ToNot(BeNil())
	Expect(result.Cachable).To(BeTrue())

	r = models.RequestDetails{
		Method:      "POST",
		Destination: "miss",
		Query: map[string][]string{
			"foo": {"bar"},
		},
		Scheme: "http",
		Body:   "body",
		Path:   "/foo",
		Headers: map[string][]string{
			"miss": {"me"},
		},
	}

	result = matching.FirstMatchStrategy(r, false, simulation, make(map[string]string))

	Expect(result.Error).ToNot(BeNil())
	Expect(result.Cachable).To(BeTrue())

	r = models.RequestDetails{
		Method:      "POST",
		Destination: "www.test.com",
		Query: map[string][]string{
			"miss": []string{""},
		},
		Scheme: "http",
		Body:   "body",
		Path:   "/foo",
		Headers: map[string][]string{
			"miss": {"me"},
		},
	}

	result = matching.FirstMatchStrategy(r, false, simulation, make(map[string]string))

	Expect(result.Error).ToNot(BeNil())
	Expect(result.Cachable).To(BeTrue())

	r = models.RequestDetails{
		Method:      "POST",
		Destination: "www.test.com",
		Query: map[string][]string{
			"foo": {"bar"},
		},
		Scheme: "http",
		Body:   "miss",
		Path:   "/foo",
		Headers: map[string][]string{
			"miss": {"me"},
		},
	}

	result = matching.FirstMatchStrategy(r, false, simulation, make(map[string]string))

	Expect(result.Error).ToNot(BeNil())
	Expect(result.Cachable).To(BeTrue())

	r = models.RequestDetails{
		Method:      "POST",
		Destination: "www.test.com",
		Query: map[string][]string{
			"foo": {"bar"},
		},
		Scheme: "http",
		Body:   "body",
		Path:   "miss",
		Headers: map[string][]string{
			"miss": {"me"},
		},
	}

	result = matching.FirstMatchStrategy(r, false, simulation, make(map[string]string))

	Expect(result.Error).ToNot(BeNil())
	Expect(result.Cachable).To(BeTrue())
}

func Test_FirstMatchStrategy_RequestMatchersShouldMatchOnStateAndNotBeCachable(t *testing.T) {
	RegisterTestingT(t)

	simulation := models.NewSimulation()

	simulation.AddRequestMatcherResponsePair(&models.RequestMatcherResponsePair{
		RequestMatcher: models.RequestMatcher{
			RequiresState: map[string]string{"key1": "value1", "key2": "value2"},
		},
		Response: testResponse,
	})

	r := models.RequestDetails{
		Body: "body",
	}

	result := matching.FirstMatchStrategy(
		r,
		false,
		simulation,
		map[string]string{"key1": "value1", "key2": "value2"})

	Expect(result.Error).To(BeNil())
	Expect(result.Cachable).To(BeFalse())
	Expect(result.Pair.Response.Body).To(Equal("request matched"))
}

func Test_FirstMatchShouldNotBeCachableIfMatchedOnEverythingApartFromStateAtLeastOnce(t *testing.T) {
	RegisterTestingT(t)

	simulation := models.NewSimulation()

	simulation.AddRequestMatcherResponsePair(&models.RequestMatcherResponsePair{
		RequestMatcher: models.RequestMatcher{
			Method: &models.RequestFieldMatchers{
				Matcher:    "exact",
				Value:      "POST",
				ExactMatch: StringToPointer("POST"),
			},
			Body: &models.RequestFieldMatchers{
				Matcher:    "exact",
				Value:      "body",
				ExactMatch: StringToPointer("body"),
			},
			Scheme: &models.RequestFieldMatchers{
				Matcher:    "exact",
				Value:      "http",
				ExactMatch: StringToPointer("http"),
			},
			Query: &models.RequestFieldMatchers{
				Matcher:    "exact",
				Value:      "foo=bar",
				ExactMatch: StringToPointer("foo=bar"),
			},
			Path: &models.RequestFieldMatchers{
				Matcher:    "exact",
				Value:      "/foo",
				ExactMatch: StringToPointer("/foo"),
			},
			Destination: &models.RequestFieldMatchers{
				Matcher:    "exact",
				Value:      "www.test.com",
				ExactMatch: StringToPointer("www.test.com"),
			},
			RequiresState: map[string]string{
				"foo": "bar",
			},
		},
		Response: testResponse,
	})

	simulation.AddRequestMatcherResponsePair(&models.RequestMatcherResponsePair{
		RequestMatcher: models.RequestMatcher{
			Method: &models.RequestFieldMatchers{
				Matcher:    "exact",
				Value:      "GET",
				ExactMatch: StringToPointer("GET"),
			},
		},
		Response: testResponse,
	})

	r := models.RequestDetails{
		Method:      "POST",
		Destination: "www.test.com",
		Query: map[string][]string{
			"foo": {"bar"},
		},
		Scheme: "http",
		Body:   "body",
		Path:   "/foo",
	}

	result := matching.FirstMatchStrategy(r, false, simulation, map[string]string{"miss": "me"})

	Expect(result.Error).ToNot(BeNil())
	Expect(result.Cachable).To(BeFalse())
}

func Test_FirstMatchShouldBeCachableIfMatchedOnEverythingApartFromStateZeroTimes(t *testing.T) {
	RegisterTestingT(t)

	simulation := models.NewSimulation()

	simulation.AddRequestMatcherResponsePair(&models.RequestMatcherResponsePair{
		RequestMatcher: models.RequestMatcher{
			Method: &models.RequestFieldMatchers{
				Matcher:    "exact",
				Value:      "POST",
				ExactMatch: StringToPointer("POST"),
			},
			Body: &models.RequestFieldMatchers{
				Matcher:    "exact",
				Value:      "body",
				ExactMatch: StringToPointer("body"),
			},
			Scheme: &models.RequestFieldMatchers{
				Matcher:    "exact",
				Value:      "http",
				ExactMatch: StringToPointer("http"),
			},
			Query: &models.RequestFieldMatchers{
				Matcher:    "exact",
				Value:      "?foo=bar",
				ExactMatch: StringToPointer("?foo=bar"),
			},
			Path: &models.RequestFieldMatchers{
				Matcher:    "exact",
				Value:      "/foo",
				ExactMatch: StringToPointer("/foo"),
			},
			Destination: &models.RequestFieldMatchers{
				Matcher:    "exact",
				Value:      "www.test.com",
				ExactMatch: StringToPointer("www.test.com"),
			},
			RequiresState: map[string]string{
				"foo": "bar",
			},
		},
		Response: testResponse,
	})

	simulation.AddRequestMatcherResponsePair(&models.RequestMatcherResponsePair{
		RequestMatcher: models.RequestMatcher{
			Method: &models.RequestFieldMatchers{
				Matcher:    "exact",
				Value:      "GET",
				ExactMatch: StringToPointer("GET"),
			},
		},
		Response: testResponse,
	})

	r := models.RequestDetails{
		Method:      "MISS",
		Destination: "www.test.com",
		Query: map[string][]string{
			"foo": {"bar"},
		},
		Scheme: "http",
		Body:   "body",
		Path:   "/foo",
	}

	result := matching.FirstMatchStrategy(r, false, simulation, map[string]string{"miss": "me"})

	Expect(result.Error).ToNot(BeNil())
	Expect(result.Cachable).To(BeTrue())

	r = models.RequestDetails{
		Method:      "POST",
		Destination: "miss",
		Query: map[string][]string{
			"foo": {"bar"},
		},
		Scheme: "http",
		Body:   "body",
		Path:   "/foo",
	}

	result = matching.FirstMatchStrategy(r, false, simulation, map[string]string{"miss": "me"})

	Expect(result.Error).ToNot(BeNil())
	Expect(result.Cachable).To(BeTrue())

	r = models.RequestDetails{
		Method:      "POST",
		Destination: "www.test.com",
		Query: map[string][]string{
			"miss": []string{""},
		},
		Scheme: "http",
		Body:   "body",
		Path:   "/foo",
	}

	result = matching.FirstMatchStrategy(r, false, simulation, map[string]string{"miss": "me"})

	Expect(result.Error).ToNot(BeNil())
	Expect(result.Cachable).To(BeTrue())

	r = models.RequestDetails{
		Method:      "POST",
		Destination: "www.test.com",
		Query: map[string][]string{
			"foo": {"bar"},
		},
		Scheme: "http",
		Body:   "miss",
		Path:   "/foo",
	}

	result = matching.FirstMatchStrategy(r, false, simulation, map[string]string{"miss": "me"})

	Expect(result.Error).ToNot(BeNil())
	Expect(result.Cachable).To(BeTrue())

	r = models.RequestDetails{
		Method:      "POST",
		Destination: "www.test.com",
		Query: map[string][]string{
			"foo": {"bar"},
		},
		Scheme: "http",
		Body:   "body",
		Path:   "miss",
	}

	result = matching.FirstMatchStrategy(r, false, simulation, map[string]string{"miss": "me"})

	Expect(result.Error).ToNot(BeNil())
	Expect(result.Cachable).To(BeTrue())
}
