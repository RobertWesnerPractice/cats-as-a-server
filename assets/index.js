window.addEventListener('load', () => {
    const load = image => {
        document.querySelector('#view').src = `/api/get/${image}`;
    };

    fetch('/api/list')
        .then(r => r.json())
        .then(images => {
            const list = document.querySelector('#list');
            list.textContent = '';

            images.forEach(image => {
                const element = document.createElement('li');
                element.textContent = image;
                element.addEventListener('click', () => load(image));
                list.append(element);
            });
        });
});
