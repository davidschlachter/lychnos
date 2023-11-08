import React from 'react';
import AmountLeft from './AmountLeft';
import FillBar from './FillBar.js'
import Header from './Header.js'
import Spinner from './Spinner.js'
import Box from '@mui/material/Box';
import Paper from '@mui/material/Paper';
import Table from '@mui/material/Table';
import TableBody from '@mui/material/TableBody';
import TableCell from '@mui/material/TableCell';
import TableContainer from '@mui/material/TableContainer';
import TableHead from '@mui/material/TableHead';
import TableRow from '@mui/material/TableRow';
import Typography from '@mui/material/Typography';
import { Link } from 'react-router-dom';

class CategoryDetail extends React.Component {
    constructor(props) {
        super(props);
        this.state = {
            error: null,
            isLoaded: false,
            details: []
        };
    }

    componentDidMount() {
        fetch("/api/reports/categorysummary/" + String(this.props.categoryId))
            .then(res => res.json())
            .then(
                (result) => {
                    this.setState({
                        isLoaded: true,
                        details: result
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
        const { error, isLoaded, details } = this.state;
        if (error) {
            return <div>Error: {error.message}</div>;
        } else if (!isLoaded) {
            return (
                <>
                    <Header back_visibility="visible" title="Category details"></Header>
                    <Spinner />
                </>
            );
        } else if (typeof details !== 'undefined' && "error" in details) {
            return <div>Error: {details.error}</div>;
        } else {
            let timeSpent = 0
            let start = new Date(details[0].start)
            let end = new Date(details[0].end)
            let now = new Date()
            timeSpent = (Math.abs(now - start) / Math.abs(end - start)) * 100
            let totalSpent = 0
            details[0].totals.map(item => (
                totalSpent += Math.round(parseFloat(item.earned) + parseFloat(item.spent))
            ));
            const actualLabel = totalSpent > 0 ? 'Earned so far' : 'Spent so far';
            return (
                <>
                    <Header back_visibility="visible" title="Category details"></Header>
                    <Box sx={{ p: 2, pb: 8 }}>
                        <Typography variant="h6" component="div" align="center" gutterBottom>
                            {details[0].name}
                        </Typography>
                        <Typography variant="subtitle1" component="div" align="center" gutterBottom>
                            Budgeted: {details[0].amount}, {actualLabel}: {totalSpent}<br />
                            Left per month: <AmountLeft amount={details[0].amount} sum={totalSpent} timeSpent={timeSpent} />
                        </Typography>
                        <div style={{ "width": "100%;" }}>
                            <FillBar amount={details[0].amount} sum={details[0].sum} now={timeSpent}></FillBar>
                        </div>
                        <br />
                        <TableContainer component={Paper}>
                            <Table style={{ "width": "100%" }}>
                                <TableHead>
                                    <TableRow>
                                        <TableCell>Month</TableCell>
                                        <TableCell>Sum</TableCell>
                                    </TableRow>
                                </TableHead>
                                <TableBody>
                                    {details[0].totals.map(item => (
                                        <TableRow key={item.start} to={"/txns/" + details[0].id + "?start=" + encodeURIComponent(new Date(item.start).toLocaleDateString('en-CA')) + "&end=" + encodeURIComponent(new Date(item.end).toLocaleDateString('en-CA'))} component={Link}>
                                            <TableCell>{Intl.DateTimeFormat('en', { month: 'long' }).format(new Date(item.start.replace("Z", "")))}</TableCell>
                                            <TableCell>{Math.round(parseFloat(item.earned) + parseFloat(item.spent))}</TableCell>
                                        </TableRow>
                                    ))}
                                </TableBody>
                            </Table>
                        </TableContainer>
                    </Box>
                </>
            );
        }
    }
}

export default CategoryDetail;