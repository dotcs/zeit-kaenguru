const rootEl = document.getElementById("root");
const url = "https://gist.githubusercontent.com/dotcs/fd3d8440ca4e338cd0185caadcd1a009/raw/comics.json";
// const url = "https://gist.githubusercontent.com/dotcs/3c63bae05e85888621a08edb89fd87da/raw/comics.json";  // dev

/**
 * Template for a comic renderd as a card.
 */
const comicCardTpl = entry => {
    const ratio = entry.img.ratio;
    const width = Math.max(window.innerWidth, 1024);
    const height = Math.floor(width / ratio);
    const url = entry.img.src.split('/');
    url.pop();
    const fittingSizeUrl = [...url, `original__${width}x${height}__ffffff`].join('/')

    return `\
        <div class="card" id="${entry.id}">
            <div class="card__img-wrapper">
                <img data-src="${fittingSizeUrl}" width="${width}" height="${height}" />
            </div>
            <div class="card__meta">
                <a href="#${entry.id}"><h2 class="card__headline">#${entry.id} ${entry.title}</h2></a>
                <p>Released: ${entry.date.substr('0', '2020-01-01'.length)}</p>
            </div>
        </div>`;
}

/** 
 * Lazy loads images when they enter the scroll area.
 * Inspired by: https://webdevtrick.com/lazy-load-images/
 */
function setupImgLazyLoad() {
    const elements = document.querySelectorAll('img[data-src]');
    let index = 0;
    const lazyLoad = function () {
        if (index >= elements.length) return;
        const item = elements[index];
        const parent = item.parentElement;
        if ((this.scrollY + this.innerHeight) > parent.offsetTop) {
            const src = item.getAttribute("data-src");
            item.src = src;
            item.addEventListener('load', function () {
                item.removeAttribute('data-src');
                item.hasAttribute('width') && item.removeAttribute('width');
                item.hasAttribute('height') && item.removeAttribute('height');
                parent.classList.add("loaded");
            });
            index++;
            lazyLoad();
        }
    };
    const init = function () {
        window.addEventListener('scroll', lazyLoad);
        lazyLoad();
    };
    return init();
}

fetch(url)
    .then(res => res.json())
    .then(items => {
        items.reverse(); 

        const html = `<main class="main">${items.map(comicCardTpl).join('')}</main>`;
        rootEl.innerHTML = html;

        setupImgLazyLoad();
    });
