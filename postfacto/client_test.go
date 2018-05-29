package postfacto

import (
	"errors"
	"fmt"
	"net/http"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/onsi/gomega/ghttp"
)

var _ = Describe("Client", func() {
	Describe("Get", func() {
		var (
			server *Server
			team   string
			client *Client

			handlers []http.HandlerFunc

			actualErr   error
			actionItems []ActionItem
		)

		BeforeEach(func() {
			server = NewServer()
			team = "some-team"

			client = NewClient(server.URL())
		})

		JustBeforeEach(func() {
			server.AppendHandlers(
				CombineHandlers(handlers...),
			)

			actionItems, actualErr = client.Get(team)
		})

		Context("when the request succeeds", func() {
			BeforeEach(func() {
				handlers = []http.HandlerFunc{
					VerifyRequest("GET", fmt.Sprintf("/retros/%s", team)),
					RespondWith(http.StatusOK, `{
						"retro":
							{
								"action_items":[
									{"id":8962, "description":"(PersonA) Leverage retro bot to post AI to slack", "done":false, "created_at":"2018-04-06T23:30:30.596Z"},
									{"id":8971, "description":"(PersonB) some-retro-message", "done":false, "created_at":"2018-04-06T23:43:40.526Z"}
								]
							}
						}`,
					),
				}
			})

			It("returns the array of action items", func() {
				Expect(server.ReceivedRequests()).Should(HaveLen(1))
				Expect(actualErr).ToNot(HaveOccurred())
				Expect(actionItems).To(ConsistOf(
					ActionItem{Description: "Leverage retro bot to post AI to slack", Member: "PersonA"},
					ActionItem{Description: "some-retro-message", Member: "PersonB"},
				))
			})
		})

		Context("when the request fails", func() {
			BeforeEach(func() {
				handlers = []http.HandlerFunc{
					VerifyRequest("GET", fmt.Sprintf("/retros/%s", team)),
					RespondWith(http.StatusTeapot, nil),
				}
			})

			It("raises an error", func() {
				Expect(server.ReceivedRequests()).Should(HaveLen(1))
				Expect(actualErr).To(MatchError(errors.New("unexpected HTTP status: 418")))
			})
		})
	})
})
