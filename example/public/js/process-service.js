import  "./mithril.min.js"

export const Process = {
    list: [],
    statusList: ["open", "rejected", "in_progress", "completed"],
    currentUUID: "",
    current: {
        firstName: "",
        lastName: "",
        email: "",
    },
    currentProcess: {
        uuid: "",
        statuses: [],
        currentStatus: "",
        current_status: {},
        payload: {}
    },
    newStatus: {
        uuid: "",
        statusName: "",
        payload: {}
    },
    error: "",
    getProcess: function(uuid) {
        return m.request({
            method: "GET",            
            url: `/api/v1/process/requests/${uuid}`,
            withCredentials: true
        }).
        then(function(res){
            if (res.length > 0) {
                console.log(res);
                Process.currentProcess = res[0]
                Process.currentProcess.currentStatus = Process.currentProcess.statuses[Process.currentProcess.statuses.length - 1]
            }            
            
        })
    },

    setStatus: function() {
        var uuid = Process.currentProcess.uuid
        var status = Process.newStatus.statusName
        var payload = Process.newStatus.payload
        if (typeof payload != "object") {
            payload = payload.replace(/\s/g, "")
            // payload = `{${payload}}`
            var data = JSON.parse(payload)
            payload = {
                payload: {
                    data : data
                }
            }
        }        

        
        console.log("PATCH",payload);

        return m.request({
            method: "PATCH",            
            url: `/api/v1/process/requests/${uuid}/assign/${status}`,
            body: payload,
            withCredentials: true
        }).
        then(function(res){
            Process.error = ""
            Process.getProcess(Process.currentProcess.uuid)
        }).
        catch(function(e) {
            console.error(e.response.message)
            Process.error = e.response.message
        })
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

        var payload = {
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
            body: payload,
            withCredentials: true
        }).
        then(function(res){
            Process.currentUUID = res.data
        })
    }
}