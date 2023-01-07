for (let index = 0; index < 6; index++) {   
    let card = document.getElementById('card-count').innerHTML += `
    <div class="col" >
        <div class="card p-2">
            <img src="/public/assets/img/pict1.jpg" class="card-img-top rounded" alt="">
            <div class="card-body">
                <a target="_blank" class="card-title h5" href="/project/${index}">Dumbways Mobile App</a>
                <div class="card-duration mb-4 text-secondary">Durasi : 3 Bulan</div>
                <div class="card-text">Lorem ipsum dolor, sit amet consectetur adipisicing elit. Sunt, officia.</div>
                <div class="card-icon mt-2 fs-3">
                    <i class="fa-brands fa-google-play"></i>
                    <i class="fa-brands fa-android"></i>
                    <i class="fa-brands fa-java"></i>
                </div>
                <div class="btn-project mt-2">
                    <a href="#" class="btn btn-dark">Edit</a>
                    <a href="#" class="btn btn-dark">Delete  </a>
                </div>
            </div>
        </div>
    </div>
    `
}