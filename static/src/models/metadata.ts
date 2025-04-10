
// import { request } from '@umijs/max';
// import { useEffect, useState } from 'react';

// export default function Page() {
//     const [brands, setBrands] = useState({});
//     const [colors, setColors] = useState({});
//     const [types, setTypes] = useState({});
//     const [loading, setLoading] = useState(true);

//     useEffect(() => {
//         request("/api/v1/meta-data").then((res) => {
//             console.log(res);
//             setBrands(res.data.brands);
//             setColors(res.data.colors);
//             setTypes(res.data.types);
//         }).catch(() => {

//         }).finally(() => {
//             setLoading(false);
//         });
//     }, [loading]);

//     return {
//         brands,
//         colors,
//         types,
//         loading,
//     };
// };