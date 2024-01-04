import "./mithril.min.js"

export const NavigationMenu = {
  view: function () {
    return m("header",
      m("nav",
        [
          m("ul",
            m("li",
              m("strong",
                "Registration Process"
              )
            )
          ),
          m("ul",
            [
              m("li",
                m("a[href='#!/list']",
                  "List"
                )
              ),
              m("li",
                m("a[href='#!/change']",
                  "Change Status"
                )
              ),
              m("li",
                m("a[href='#!/create']",
                  "Create New"
                )
              ),
              // m("li", 
              //   m("a[href='#'][role='button']", 
              //     "Button"
              //   )
              // )
            ]
          )
        ]
      )
    )
  }
}