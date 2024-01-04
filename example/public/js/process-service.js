import  "./mithril.min.js"

export const Process = {
    list: [],
    currentUUID: "",
    current: {
        firstName: "",
        lastName: "",
        email: "",
    },
    loadList: function() {
        return m.request({
            method: "GET",
            url: `/api/v1/process/requests/list`,
            withCredentials: true
        }).
        then(function(res){
            Process.list = res.data
        })
    },
    create: function() {

        var req = {
            code: "requests",
            current_status: {
                name: "open"
            },
            payload: {
                firstName: Process.current.firstName,
                lastName: Process.current.lastName,
                email: Process.current.email
            }
        }
        return m.request({
            method: "POST",
            url: `/api/v1/process`,
            body: req,
            withCredentials: true
        }).
        then(function(res){
            Process.currentUUID = res.data
        })
    }
}