import { useState, useEffect } from "react";

const queryString = (params) =>
    Object.keys(params)
        .map((key) => `${key}=${params[key]}`)
        .join("&");

const createUrl = (url, queryOptions) => {
    return url + "?" + queryString(queryOptions);
};

export default (url, options = { body: {}, query: {} }) => {
    const [data, setData] = useState({
        response: null,
        error: false,
        loading: true,
    });

    useEffect(() => {
        setData({ ...data, error: null, loading: true });
        fetch(createUrl(url, options.query), {
            method: options.method || "GET",
            headers: {
                "Content-Type": "application/json",
            },
            body: options.method !== "GET" && JSON.stringify(options.body),
        })
            .then(async (response) => {
                const data = await response.json();
                setData({
                    response: data,
                    error: !response.ok,
                    loading: false,
                });
            })
            .catch((error) => {
                //fetch throws an error only on network failure or if anything prevented the request from completing
                setData({
                    response: { status: "network_failure" },
                    error: true,
                    loading: false,
                });
            });
    }, [url, JSON.stringify(options)]);

    return data;
};


