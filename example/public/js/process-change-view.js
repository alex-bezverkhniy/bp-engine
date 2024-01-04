import "./mithril.min.js"
import { Process } from './process-service.js'
import { NavigationMenu } from './navigation.js'

export const ProcessChangeForm = {
    oninit: Process.loadList,
    view: function () {
        
        return [
            m("main",
                m(NavigationMenu),
                m("h1", {
                    class: "title"
                }, "Change Process Status"),

                m("form", {
                    onsubmit: function (e) {
                        e.preventDefault()
                        Process.setStatus()
                    }
                },
                    [
                        m("div.grid",
                            [
                                m("label[for='uuid']",
                                    [
                                        " Select Process: ",
                                        m("select[id='uuid'][required]",
                                            {
                                                onchange: function (e) {
                                                    Process.newStatus.uuid = e.target.options[e.target.options.selectedIndex].value
                                                    console.log(Process.newStatus.uuid)
                                                    Process.getProcess(Process.newStatus.uuid)
                                                    Process.newStatus.statusName = ''
                                                    Process.error = ''
                                                },
                                                value: Process.newStatus.uuid
                                            },
                                            [
                                                m("option[value='']",
                                                    "Select a process…"
                                                ),
                                                Process.list.map(function (p) {
                                                    return m(`option[value='${p.uuid}']`, p.uuid)
                                                })
                                            ]
                                        )
                                    ]
                                ),
                                m("label[for='curStatus']",
                                    [
                                        " Current Status ",
                                        m("input[type='text'][id='curStatus'][name='curStatus'][placeholder='Current Status'][readonly]",
                                            {                                                
                                                value: Process.currentProcess.current_status.name
                                            }
                                        )
                                    ]
                                ),
                                m("label[for='status']",
                                    [
                                        " Select Next Status: ",
                                        m("select[id='status'][required]",
                                            {
                                                onchange: function (e) {
                                                    e.preventDefault()
                                                    Process.newStatus.statusName = e.target.options[e.target.options.selectedIndex].value
                                                    console.log(Process.newStatus.statusName)

                                                },
                                                value: Process.newStatus.statusName
                                            },
                                            [
                                                m("option[value='']",
                                                    "Select a status…"
                                                ),
                                                Process.statusList.map(function (s) {
                                                    return m(`option[value='${s}']`, s)
                                                })
                                            ]
                                        )
                                    ]
                                ),
                                m("label[for='payload']",
                                    [
                                        " Payload ",
                                        m("textarea[id='payload'][name='payload'][placeholder='JSON payload'][required]",
                                            {
                                                onchange: function(e) {Process.newStatus.payload = e.target.value},
                                                value: typeof Process.newStatus.payload == "object" ? JSON.stringify(Process.newStatus.payload, undefined, 2): Process.newStatus.payload
                                            }
                                        )
                                    ]
                                ),

                            ]
                        ),
                        m("kbd.error", {style:{display: Process.error == "" ? 'none': ''}}, Process.error),
                        m("button[type='submit']",
                            "Submit"
                        )
                    ]
                )
            )
        ]
    }
}