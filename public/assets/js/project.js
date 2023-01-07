let projects = []
function getProject(event){
    event.preventDefault()
    let name    = document.getElementById('name').value
    let sDate   = new Date(document.getElementById('sDate').value)
    let eDate   = new Date(document.getElementById('eDate').value)
    let durasi  = getDistanceTime(sDate, eDate)
    let desc    = document.getElementById('description').value
    
    var checkedVal   = []
    var checked = document.getElementsByClassName('form-check-input')
    for (let i = 0; i < checked.length; i++) {
        if (checked[i].checked) {
            checkedVal[i] = checked[i].value
        }else{
            checkedVal[i] = ""
        }
    }

    let image   = document.getElementById('image').files
    image       = URL.createObjectURL(image[0])
    
    let addProject = {
        name,
        durasi,
        desc,
        checkedVal,
        image
    }

    projects.push(addProject)
    showData()
}

function showData() {
    document.getElementById('card-count').innerHTML = ""
    for (let index = 0; index < projects.length; index++) {   
        let checkBox = ""
        for (let i = 0; i < projects[index].checkedVal.length; i++) {
            checkBox += `
            ${projects[index].checkedVal[i]}
            `
        }
        let card = document.getElementById('card-count')
        card.innerHTML += `
        <div class="col" >
            <div class="card p-2">
                <img src="${projects[index].image}" class="card-img-top rounded" alt="">
                <div class="card-body">
                    <h5 class="card-title">${projects[index].name}</h5>
                    <div class="card-duration mb-4 text-secondary">${projects[index].durasi}</div>
                    <div class="card-text">${projects[index].desc}</div>
                    <div class="card-icon mt-2 fs-3">
                        ${checkBox}
                    </div>
                    <div class="btn-project mt-2">
                        <a href="" class="btn btn-dark">Edit</a>
                        <a href="" class="btn btn-dark">Delete  </a>
                    </div>
                </div>
            </div>
        </div>
        `
    }
}

function getDistanceTime(start, end){
    let timePost    = start
    let timeNow     = end
    let distance    = timeNow - timePost
    let disMonth    = Math.floor(distance / (1000 * 60 * 60 * 24 * 30))
    let disDay      = Math.floor(distance / (1000 * 60 * 60 * 24))
    let disHour     = Math.floor(distance / (1000 * 60 * 60))
    let disMinute   = Math.floor(distance / (1000 * 60))
    let disSecond   = Math.floor(distance / (1000))
    // distance = distance / 1000
    // return `${Math.floor(distance)} Second Ago`
    if (disMonth > 0) {
        return `Durasi: ${disMonth} Bulan`
    } else if(disDay > 0){
        return `Durasi: ${disDay} Hari`
    } else if(disHour > 0){
        return `Durasi: ${disHour} Jam`
    } else if(disMinute > 0){
        return `Durasi: ${disMinute} Menit`
    } else if(disSecond > 0){
        return `Durasi: ${disSecond} Detik`
    }
}