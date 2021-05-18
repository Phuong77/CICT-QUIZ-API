console.log("Hello word");
const dropdownContainer = document.getElementById('dropdown-container');

const url = 'http://localhost:9000/api/categories';

fetch(url)
    .then(response => response.json())
    .then(result => {
        console.log(result);

        const status = result.status;

        // console.log('Data');
        // console.log(data);
        //
        // console.log('status');
        // console.log(status);

        if (status !== "success") {
            alert('Có lỗi trong quá trình lấy dữ liệu từ máy chủ');
            return;
        }

        const data = result.data;

        const parents = [];
        const children = [];

        for (const item of data) {
            if (item.parent_id === "") {
                parents.push(item);
            }
            else {
                children.push(item);
            }
        }

        console.log(parents);
        console.log(children);


        for (const parent of parents) {
            const dropDownElement = document.createElement('div');
            dropDownElement.classList.add('dropdown');

            const dropButtonElement = document.createElement('button');
            dropButtonElement.classList.add('dropbtn');
            dropButtonElement.innerHTML = parent.name;

            const dropdownContentElement = document.createElement('div');
            dropdownContentElement.classList.add('dropdown-content');

            dropDownElement.appendChild(dropButtonElement);
            dropDownElement.appendChild(dropdownContentElement);

            for (const child of children){
                if(child.parent_id==parent._id){
                    const aelement = document.createElement('a');
                    aelement.innerHTML = child.name;
                    aelement.style.cursor = 'pointer';

                    aelement.addEventListener('click', function () {
                        console.log(child.name + ' click');
                    });

                    dropdownContentElement.appendChild(aelement);
                }
            }

            dropdownContainer.appendChild(dropDownElement);
        }

        // for (let i=0;i<parents.length;i++){
        //     const parent=parents[i];
        //     const html = `
        //          <div class="dropdown">
        //             <button class="dropbtn">${parent.name}</button>
        //             <div class="dropdown-content">
        //                 <a href="#">Quốc Phòng 1</a>
        //                 <a href="#">Quốc phòng 2</a>
        //                 <a href="#">Quốc phòng 3</a>
        //                 <a href="#">Quốc phòng 4</a>
        //             </div>
        //         </div>
        //     `;
        //
        //     dropdownContainer.insertAdjacentHTML('afterbegin', html);
        // }



    });