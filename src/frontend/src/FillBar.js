import './FillBar.css';


function FillBar(props) {
    let width = props.sum / props.amount * 100
    let color = 'green'
    if (width > 100) {
        width = 100
        color = 'red'
    }
    let style = {
        'width': String(width) + "%",
        'backgroundColor': color
    }

    let nowStyle = {
        'left': String(props.now) + "%",
        'width': "5px",
        'height': "1em",
        'position': "relative",
        'backgroundColor': 'black',
        'top': "-1em"
    }

    return (
        <div className="fillBar">
            <div className="filled" style={style}></div>
            <div className="nowLine" style={nowStyle}></div>
        </div>
    );
}

export default FillBar;
