import * as React from 'react';
import './TransactionList.css';
import AccountIcon from './AccountIcon.js';
import FriendlyDate from './FriendlyDate.js';
import Spinner from './Spinner.js';
import useFetch from "./useFetch";
import Box from '@mui/material/Box';
import Header from './Header.js';
import Table from '@mui/material/Table';
import TableBody from '@mui/material/TableBody';
import TableCell from '@mui/material/TableCell';
import TableContainer from '@mui/material/TableContainer';
import TableHead from '@mui/material/TableHead';
import TableRow from '@mui/material/TableRow';
import { Link, useParams, useSearchParams } from 'react-router-dom';

export default function TransactionList(props) {
    const { categoryID } = useParams();
    let back_visibility;
    if (typeof categoryID === 'undefined') {
        back_visibility = "hidden";
    } else {
        back_visibility = "visible";
    }

    const [searchParams, setSearchParams] = useSearchParams();
    let start = searchParams.get("start");
    let end = searchParams.get("end");

    // Default date range: today - 30 days
    if (start == null) {
        start = new Date(new Date().setDate(new Date().getDate() - 30)).toLocaleDateString('en-CA');
    }
    if (end == null) {
        end = new Date().toLocaleDateString('en-CA')
    }

    const { response, error, loading } = useFetch(
        "/api/transactions",
        {
            query: {
                start: start,
                end: end
            },
        }
    );

    if (loading) {
        return (
            <>
                <Header back_visibility={back_visibility} title="Transactions"></Header>
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

    let txns;
    if (typeof categoryID !== 'undefined') {
        txns = [];
        for (const t of response) {
            if (t.attributes.transactions[0].category_id !== categoryID.toString()) {
                continue;
            }
            txns.push(t);
        }
    } else {
        txns = response;
    }

    return (
        <>
            <Header back_visibility={back_visibility} title="Transactions"></Header>
            <Box sx={{ pb: 8 }}>
                <TableContainer>
                    <Table className="txnList" style={{ "width": "100%" }}>
                        <TableHead>
                            <TableRow>
                                <TableCell>Description</TableCell>
                                <TableCell className="txnListAmount">Amount</TableCell>
                                <TableCell className="txnListFrom">From</TableCell>
                                <TableCell>To</TableCell>
                                <TableCell>Date</TableCell>
                            </TableRow>
                        </TableHead>
                        <TableBody>
                            {txns.map(item => (
                                <TableRow key={item.id} to={"/txn/" + item.id} component={Link}>
                                    <TableCell>{item.attributes.transactions[0].description}</TableCell>
                                    <TableCell>{parseFloat(item.attributes.transactions[0].amount).toLocaleString(undefined, { maximumFractionDigits: 2, minimumFractionDigits: 2 })}</TableCell>
                                    <TableCell><AccountIcon account_id={item.attributes.transactions[0].source_id} /> <span className="srcName">{item.attributes.transactions[0].source_name}</span></TableCell>
                                    <TableCell>{item.attributes.transactions[0].destination_name}</TableCell>
                                    <TableCell><FriendlyDate date={item.attributes.transactions[0].date} /></TableCell>
                                </TableRow>
                            ))}
                        </TableBody>
                    </Table>
                </TableContainer>
            </Box>
        </>
    );
}