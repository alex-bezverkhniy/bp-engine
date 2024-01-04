import "./mithril.min.js"
import { Process } from './process-service.js'
import { NavigationMenu } from './navigation.js'

export const ProcessCreateForm = {
    // oninit: Process.loadList,
    view: function () {
        return [
            m("main",
                m(NavigationMenu),
                m("h1", {
                    class: "title"
                }, "Create New Processes"),

                m("form",{
                    onsubmit: function(e) {
                        e.preventDefault()
                        Process.create()
                    }
                },
                    [
                        m("div.grid",
                            [
                                m("label[for='firstname']",
                                    [
                                        " First name ",
                                        m("input[type='text'][id='firstname'][name='firstname'][placeholder='First name'][required]",
                                        {
                                            oninput: function(e) {Process.current.firstName = e.target.value},
                                            value: Process.current.firstName
                                        })
                                    ]
                                ),
                                m("label[for='lastname']",
                                    [
                                        " Last name ",
                                        m("input[type='text'][id='lastname'][name='lastname'][placeholder='Last name'][required]",
                                        {
                                            oninput: function(e) {Process.current.lastName = e.target.value},
                                            value: Process.current.lastName
                                        })
                                    ]
                                )
                            ]
                        ),
                        m("label[for='email']",
                            "Email address"
                        ),
                        m("input[type='email'][id='email'][name='email'][placeholder='Email address'][required]",
                        {
                            oninput: function(e) {Process.current.email = e.target.value},
                            value: Process.current.email
                        }),
                        m("small",
                            "We'll never share your email with anyone else."
                        ),
                        m("button[type='submit']",
                            "Submit"
                        )
                    ]
                )
            )
        ]
    }
}