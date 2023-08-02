function listGet(id) {
  return axios({
    method: "GET",
    url: `/api/list/${id}`
  })
}

function listList() {
  return axios({
    method: "GET",
    url: `/api/list`
  })
}

function listPost(list) {
  let submitButton = document.querySelector("#post-form-submit")
  submitButton.disabled = true
  let name = document.querySelector('#post-form-name').value
  let description = document.querySelector('#post-form-description').value
  return axios({
    method: "POST",
    url: `/api/list`,
    data: {
      'name': name,
      'description': description
    }
  }).then(resp => {
    submitButton.disabled = false
    window.location.reload(false)
  }).catch(err => {
    submitButton.disabled = false
  })
}

function listPut(id, list) {
  return axios({
    method: "PUT",
    url: `/api/list/${id}`,
    data: list
  })
}

function listDelete(id) {
  return axios({
    method: "DELETE",
    url: `/api/list/${id}`
  })
}

function listItemGet(listid, id) {
  return axios({
    method: "GET",
    url: `/api/list/${listid}/item/${id}`
  })
}

function listItemList(listid) {
  return axios({
    method: "GET",
    url: `/api/list/${listid}/item`
  })
}

function listItemPost(listid, item) {
  let submitButton = document.querySelector("#post-form-submit")
  submitButton.disabled = true
  let name = document.querySelector('#post-form-name').value
  let description = document.querySelector('#post-form-description').value
  const urlSearchParams = new URLSearchParams(window.location.search)
  const params = Object.fromEntries(urlSearchParams.entries())
  return axios({
    method: "POST",
    url: `/api/list/${params.id}/item`,
    data: {
      'name': name,
      'description': description
    }
  }).then(resp => {
    submitButton.disabled = false
    window.location.reload(false)
  }).catch(err => {
    submitButton.disabled = false
  })
}

function listItemPut(listid, id, item) {
  return axios({
    method: "PUT",
    url: `/api/list/${listid}/item/${id}`,
    data: list
  })
}

function listItemDelete(listid, id) {
  return axios({
    method: "DELETE",
    url: `/api/list/${listid}/item/${id}`
  })
}

function listItemDeleteAll(listid) {
  return axios({
    method: "DELETE",
    url: `/api/list/${listid}/item`
  })
}

function todoListsPage() {
  let form = document.getElementById("form")
  form.onsubmit = function(event) {
    event.preventDefault()
    return false
  }
  const listsList = document.querySelector("#lists-list table")
  const template = document.querySelector("#list-item-template")
  listList().then(resp => {
    if (resp.data === null) {
      return
    }
    resp.data.forEach(list => {
      let clone = template.content.firstElementChild.cloneNode(true)
      let name = clone.querySelector(".list-item-name a")
      name.textContent = list.name
      name.href = `/list.html?id=${list.id}`
      let description = clone.querySelector(".list-item-description p")
      description.textContent = list.description
      let deleteButton = clone.querySelector(".list-item-delete button")
      deleteButton.onclick = () => {
        listDelete(list.id).then(() => {
          window.location.reload(false)
        })
      }
      listsList.append(clone)
    }).catch(err => {
      console.log({err})
    })
  })
}

function todoListPage() {
  let form = document.getElementById("form")
  form.onsubmit = function(event) {
    event.preventDefault()
    return false
  }
  const urlSearchParams = new URLSearchParams(window.location.search)
  const params = Object.fromEntries(urlSearchParams.entries())
  listGet(params.id).then(resp => {
    if (resp.data === null) {
      return
    }
    const list = resp.data
    const listNameTitle = document.querySelector("h1#list-name")
    listNameTitle.textContent = list.name
  })
  const listsList = document.querySelector("#lists-list table")
  const template = document.querySelector("#list-item-template")
  listItemList(params.id).then(resp => {
    if (resp.data === null) {
      return
    }
    resp.data.forEach(list => {
      let clone = template.content.firstElementChild.cloneNode(true)
      let name = clone.querySelector(".list-item-name p")
      name.textContent = list.name
      let description = clone.querySelector(".list-item-description p")
      description.textContent = list.description
      let deleteButton = clone.querySelector(".list-item-delete button")
      deleteButton.onclick = () => {
        listItemDelete(params.id, list.id).then(() => {
          window.location.reload(false)
        })
      }
      listsList.append(clone)
    }).catch(err => {
      console.log({err})
    })
  })
}

document.addEventListener('DOMContentLoaded',() => {
  let page = document.querySelector("meta[name=todolist-page]")?.content
  switch (page) {
    case "home":
      todoListsPage()
      break;
    case "list":
      todoListPage()
    default:
      console.log("meta[name=todolist-page] needs to be filled out or set is not routed")
  }
})
