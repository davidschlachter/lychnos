import * as React from 'react';
import './TransactionList.css';
import AccountIcon from './AccountIcon.js';
import Spinner from './Spinner.js';
import FriendlyDate from './FriendlyDate.js';
import Table from '@mui/material/Table';
import TableBody from '@mui/material/TableBody';
import TableCell from '@mui/material/TableCell';
import TableContainer from '@mui/material/TableContainer';
import TableHead from '@mui/material/TableHead';
import TableRow from '@mui/material/TableRow';
import Paper from '@mui/material/Paper';
import { Link } from 'react-router-dom';
import useFetch from "./useFetch";

export default function TransactionList() {
    const { response, error, loading } = useFetch(
        "/api/transactions",
        {
            query: {
                page: 1,
                pageSize: 100,
            },
        }
    );

    if (loading) {
        return <Spinner />;
    }
    if (error) {
        return <div className="error">{JSON.stringify(error)}</div>;
    }

    return (
        <TableContainer component={Paper}>
            <Table style={{ "width": "100%" }}>
                <TableHead>
                    <TableRow>
                        <TableCell>Description</TableCell>
                        <TableCell>Amount</TableCell>
                        <TableCell>From</TableCell>
                        <TableCell>To</TableCell>
                        <TableCell>Date</TableCell>
                    </TableRow>
                </TableHead>
                <TableBody>
                    {response.map(item => (
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
    );
}