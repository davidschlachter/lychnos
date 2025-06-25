import { useState, useEffect } from "react";

const queryString = (params) =>
    Object.keys(params)
        .map((key) => `${key}=${params[key]}`)
        .join("&");

function createUrl(url, queryOptions) {
    if (Object.keys(queryOptions).length > 0) {
        return url + "?" + queryString(queryOptions);
    } else {
        return url;
    }

}

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
            .then(res => {
                if (!res.ok) {
                    return Promise.reject(res);
                }
                return res.json();
            })
            .then(result => {
                setData({
                    response: result,
                    error: false,
                    loading: false,
                });
            })
            .catch(error => {
                var error_message = `Status ${error.status} (${error.statusText})`;
                if (typeof error.text === "function") {
                    error.text().then(textError => {
                        error_message += `; Message: ${textError}`;
                        setData({
                            response: null,
                            error: error_message,
                            loading: false,
                        });
                    }).catch(genericError => {
                        setData({
                            response: null,
                            error: `Parsing error message: ${genericError}`,
                            loading: false,
                        });
                    });
                } else {
                    setData({
                        response: null,
                        error: error_message += `; Fetch error: ${error}`,
                        loading: false,
                    });
                }
            });
    }, [url, JSON.stringify(options)]);

    return data;
};


