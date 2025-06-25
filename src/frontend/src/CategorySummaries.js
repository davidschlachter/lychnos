import React from 'react';
import './CategorySummaries.css';
import AmountLeft from './AmountLeft';
import FillBar from './FillBar.js'
import Spinner from './Spinner.js'
import Box from '@mui/material/Box';
import Header from './Header.js';
import Paper from '@mui/material/Paper';
import Table from '@mui/material/Table';
import TableBody from '@mui/material/TableBody';
import TableCell from '@mui/material/TableCell';
import TableContainer from '@mui/material/TableContainer';
import TableHead from '@mui/material/TableHead';
import TableRow from '@mui/material/TableRow';
import { Link } from 'react-router-dom';

class CategorySummaries extends React.Component {
    constructor(props) {
        super(props);
        this.state = {
            error: null,
            isLoaded: false,
            summaries_data: []
        };
    }

    componentDidMount() {
        fetch("/api/reports/categorysummary/")
            .then(res => {
                if (!res.ok) {
                    return Promise.reject(res);
                }
                return res.json();
            })
            .then(result => {
                this.setState({
                    isLoaded: true,
                    summaries_data: result
                }, () => {
                    this.handleScrollPosition();
                });
            })
            .catch(error => {
                var error_message = `Status ${error.status} (${error.statusText})`;
                if (typeof error.text === "function") {
                    error.text().then(textError => {
                        error_message += `; Message: ${textError}`;
                        this.setState({
                            isLoaded: true,
                            error: error_message,
                        });
                    }).catch(genericError => {
                        this.setState({
                            isLoaded: true,
                            error: `Parsing error message: ${genericError}`,
                        });
                    });
                } else {
                    this.setState({
                        isLoaded: true,
                        error: error_message += `; Fetch error: ${error}`,
                    });
                }
            });
    }

    // Remember scroll position.
    handleScrollPosition = () => {
        const scrollPosition = sessionStorage.getItem("scrollPosition");
        if (scrollPosition) {
            window.scrollTo(0, parseInt(scrollPosition));
            sessionStorage.removeItem("scrollPosition");
        }
    }
    handleClick = () => {
        sessionStorage.setItem("scrollPosition", window.pageYOffset);
    }

    render() {
        const { error, isLoaded, summaries_data } = this.state;
        if (error) {
            return <div>Error: {error}</div>;
        } else if (!isLoaded) {
            return (
                <>
                    <Header back_visibility="hidden" title="Category summaries" budgetedit={true}></Header>
                    <Spinner />
                </>
            );
        } else if (typeof summaries_data === 'object' &&
            !Array.isArray(summaries_data) &&
            summaries_data !== null && "error" in summaries_data) {
            return (
                <>
                    <Header back_visibility="hidden" title="Category summaries" budgetedit={true}></Header>
                    <Box sx={{ p: 2, pb: 8 }}>
                        <div>Error: {summaries_data.error}</div>
                        <div>If no budget could be found, please <a href="budget/">create a budget for this year</a>.</div>
                    </Box>
                </>
            );
        } else if (summaries_data === null) {
            return (
                <>
                    <Header back_visibility="hidden" title="Category summaries" budgetedit={true}></Header>
                    <Box sx={{ p: 2, pb: 8 }}>
                        <div>Error: no category summary data could be fetched.</div>
                    </Box>
                </>
            );
        }
        else {
            let timeSpent = 0
            let start = new Date(summaries_data[0].start)
            let end = new Date(summaries_data[0].end)
            let now = new Date()
            timeSpent = (Math.abs(now - start) / Math.abs(end - start)) * 100

            const summaries = [...summaries_data].sort((a, b) => {
                // Sort by whatever is greater for a given category: the amount
                // budgeted or the amount actually spent.
                let aKey = Math.abs(+a.amount) > Math.abs(+a.sum) ? +a.amount : +a.sum
                let bKey = Math.abs(+b.amount) > Math.abs(+b.sum) ? +b.amount : +b.sum
                return aKey > bKey ? 1 : -1
            });

            return (
                <>
                    <Header back_visibility="hidden" title="Category summaries" budgetedit={true}></Header>
                    <Box sx={{ p: 2, pb: 8 }}>
                        <TableContainer component={Paper} style={{ "overflow-y": "hidden" }}>
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
                                        <TableRow key={item.category_budget_id} to={"/categorydetail/" + item.category_budget_id} onClick={this.handleClick} component={Link}>
                                            <TableCell>{item.name}</TableCell>
                                            <TableCell width="99%"><FillBar amount={item.amount} sum={item.sum} now={timeSpent}></FillBar></TableCell>
                                            <TableCell style={{
                                                'fontFamily': "monospace",
                                                'textAlign': "right",
                                                'fontSize': "110%",
                                                'fontWeight': "600"
                                            }}><AmountLeft amount={item.amount} sum={item.sum} timeSpent={timeSpent}></AmountLeft></TableCell>
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

export default CategorySummaries;
