package main

import (
	"github.com/dalloriam/graphql-tools/analyzer"
)

const x = `query() {
}`

const r = `query() {
	Me {
		Organizations {
			List(Size: 10) {
				Hits {
					ID
					Name
					SalesChannels {
						List(Size: 10) {
							Hits {
								ID
								Organization {
									ID
									SalesChannels {
										List(Size: 10) {
											Hits {
												ID
											}
										}
									}
								}
							}
						}
					}
				}
				Total
			}
		}
	}
}
`

func main() {

	schema, err := analyzer.LoadSchema("./schema")
	if err != nil {
		panic(err)
	}
	if err := ValidateQuery(r, schema); err != nil {
		panic(err)
	}
}
