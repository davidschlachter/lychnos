import React from 'react';
import FriendlyDate from './FriendlyDate.js';
import Header from './Header.js';
import Spinner from "./Spinner.js";
import useFetch from "./useFetch";
import Box from '@mui/material/Box';
import Typography from '@mui/material/Typography';
import { useParams } from "react-router-dom";

export default function TransactionDetail() {
    const { txnID } = useParams();
    const { response, error, loading } = useFetch(
        "/api/transactions/" + txnID,
        {
            query: {},
        }
    );

    if (loading) {
        return (
            <>
                <Header back_visibility="visible" back_location="/txns" title="Transaction details"></Header>
                <Spinner />
            </>
        );
    }
    if (error) {
        return <div className="error">{JSON.stringify(error)}</div>;
    }
    if (typeof response !== 'undefined' && "error" in response) {
        return <div className="error">Error: {response.error}</div>;
    }

    return (
        <>
            <Header back_visibility="visible" back_location="/txns" title="Transaction details"></Header>
            <Box sx={{ p: 2, mb: 6 }}>
                <Typography variant="h6" component="div" align="center" gutterBottom>
                    {response.attributes.transactions[0].description}
                </Typography>
                <Typography variant="subtitle1" component="div" align="left" gutterBottom>
                    From: {response.attributes.transactions[0].source_name} <br />
                    To: {response.attributes.transactions[0].destination_name} <br />
                    Date: <FriendlyDate date={response.attributes.transactions[0].date} /> <br />
                    Amount: ${parseFloat(response.attributes.transactions[0].amount).toLocaleString(undefined, { maximumFractionDigits: 2, minimumFractionDigits: 2 })} <br />
                    Category: {response.attributes.transactions[0].category_name}
                </Typography>
            </Box>
        </>
    );
}

