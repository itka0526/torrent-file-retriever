const arr = [
    // {
    //     path: "./downloads/",
    //     name: "downloads",
    //     size: 192,
    //     modified_date: "2023-11-08T06:01:48.050062+03:00",
    //     is_directory: true,
    // },
    {
        path: "downloads/photo_5309871435755409443_y.jpg",
        name: "photo_5309871435755409443_y.jpg",
        size: 177806,
        modified_date: "2023-11-08T05:54:55.66388154+03:00",
        is_directory: false,
    },
    {
        path: "downloads/test",
        name: "test",
        size: 96,
        modified_date: "2023-11-08T05:57:06.348308136+03:00",
        is_directory: true,
    },
    {
        path: "downloads/test/test1",
        name: "test1",
        size: 96,
        modified_date: "2023-11-08T06:02:10.776713467+03:00",
        is_directory: true,
    },
    {
        path: "downloads/test/test1/name.txt",
        name: "name.txt",
        size: 5,
        modified_date: "2023-11-08T01:52:07.49954395+03:00",
        is_directory: false,
    },
    {
        path: "downloads/wolf.png",
        name: "wolf.png",
        size: 76116,
        modified_date: "2023-11-06T00:57:47.945863875+03:00",
        is_directory: false,
    },
    {
        path: "downloads/МатАн3_ЛК9.pdf",
        name: "МатАн3_ЛК9.pdf",
        size: 658108,
        modified_date: "2023-11-07T02:50:24.191303674+03:00",
        is_directory: false,
    },
];

const result = {};

function s() {
    for (const item of arr) {
        const paths = item.path.split("/");

        let curr = result;
        for (let path of paths) {
            if (curr[path] && curr[path].children) {
                curr[path].children.push(item);
                curr = curr[path];
            } else if (!curr.hasOwnProperty(path)) {
                curr[path] = { ...item, children: [] };
                curr = curr[path];
            }
        }
    }
}

s();

// const item = arr[4];
// const paths = item.path.replace("/" + item.name, "").split("/");
// console.log(item.path, paths);

// let curr = result;
// for (let path of paths) {
//     if (!curr.hasOwnProperty(path)) {
//         curr[path] = {};
//     }
//     curr = curr[path];
// }

console.log(JSON.stringify(result, null, 4));
