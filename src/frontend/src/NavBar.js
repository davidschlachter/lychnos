import * as React from 'react';
import Paper from '@mui/material/Paper';
import BottomNavigation from '@mui/material/BottomNavigation';
import BottomNavigationAction from '@mui/material/BottomNavigationAction';
import AlignHorizontalLeftIcon from '@mui/icons-material/AlignHorizontalLeft';
import AddBoxIcon from '@mui/icons-material/AddBox';
import ListAltIcon from '@mui/icons-material/ListAlt';
import {
    Link,
    matchPath,
    useLocation,
} from 'react-router-dom';

export default function NavBar() {
    const routeMatch = useRouteMatch(['/newTxn', '/txns', '/']);
    const currentPage = routeMatch?.pattern?.path;

    return (
        <Paper sx={{ position: 'fixed', bottom: 0, left: 0, right: 0 }} elevation={3} style={{ zIndex: 3 }}>
            <BottomNavigation
                value={currentPage}
                showLabels
                sx={{ paddingBottom: 'env(safe-area-inset-right)' }}
            >
                <BottomNavigationAction label="Summary" icon={<AlignHorizontalLeftIcon />} value="/" to="/" component={Link} />
                <BottomNavigationAction label="New Txn" icon={<AddBoxIcon />} value="/newTxn" to="/newTxn" component={Link} />
                <BottomNavigationAction label="Txns" icon={<ListAltIcon />} value="/txns" to="/txns" component={Link} />
            </BottomNavigation >
        </Paper >
    );
}

function useRouteMatch(patterns) {
    const { pathname } = useLocation();

    for (let i = 0; i < patterns.length; i += 1) {
        const pattern = patterns[i];
        const possibleMatch = matchPath(pattern, pathname);
        if (possibleMatch !== null) {
            return possibleMatch;
        }
    }

    return null;
}