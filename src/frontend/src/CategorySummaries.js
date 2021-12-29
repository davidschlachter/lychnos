import React from 'react';
import './CategorySummaries.css';
import FillBar from './FillBar.js'
import AmountLeft from './AmountLeft';
import Table from '@mui/material/Table';
import TableBody from '@mui/material/TableBody';
import TableCell from '@mui/material/TableCell';
import TableContainer from '@mui/material/TableContainer';
import TableHead from '@mui/material/TableHead';
import TableRow from '@mui/material/TableRow';
import Paper from '@mui/material/Paper';

class CategorySummaries extends React.Component {
    constructor(props) {
        super(props);
        this.state = {
            error: null,
            isLoaded: false,
            summaries: [],
            budgets: []
        };
    }

    componentDidMount() {
        fetch("/api/budgets/")
            .then(res => res.json())
            .then(
                (result) => {
                    this.setState({
                        budgets: result
                    });
                },
                (error) => {
                    this.setState({
                        isLoaded: true,
                        error
                    });
                }
            )
        fetch("/api/reports/categorysummary/?budget=1")
            .then(res => res.json())
            .then(
                (result) => {
                    this.setState({
                        isLoaded: true,
                        summaries: result
                    });
                },
                (error) => {
                    this.setState({
                        isLoaded: true,
                        error
                    });
                }
            )
    }

    render() {
        const { error, isLoaded, summaries, budgets } = this.state;
        if (error) {
            return <div>Error: {error.message}</div>;
        } else if (!isLoaded) {
            return <div>Loading...</div>;
        } else {
            let timeSpent = 0
            if (budgets.length > 1) {
                let start = new Date(budgets[0].start)
                let end = new Date(budgets[0].end)
                let now = new Date()
                timeSpent = (Math.abs(now - start) / Math.abs(end - start)) * 100
            }
            return (
                <TableContainer component={Paper}>
                    <Table style={{ "width": "100%" }}>
                        <TableHead>
                            <TableRow>
                                <TableCell>Category</TableCell>
                                <TableCell>Progress</TableCell>
                                <TableCell>Left&nbsp;per month</TableCell>
                            </TableRow>
                        </TableHead>
                        <TableBody>
                            {summaries.map(item => (
                                <TableRow key={item.id}>
                                    <TableCell>{item.name}</TableCell>
                                    <TableCell width="99%"><FillBar amount={item.amount} sum={item.sum} now={timeSpent}></FillBar></TableCell>
                                    <AmountLeft amount={item.amount} sum={item.sum} timeSpent={timeSpent}></AmountLeft>
                                </TableRow>
                            ))}
                        </TableBody>
                    </Table>
                </TableContainer>
            );
        }
    }
}

export default CategorySummaries;
