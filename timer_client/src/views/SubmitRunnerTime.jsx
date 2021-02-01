import React, { Component } from 'react';
import { DataGrid } from '@material-ui/data-grid';
import Button from '@material-ui/core/Button';

class SubmitRunnerTime extends Component {
    constructor(props) {
        super(props);
        this.state = {
            times: [],
            chipNumber: "",
            fullName: "",
            task: "",
            columns: [
                { field: 'id', headerName: 'Chip number', width: 200 },
                { field: 'startNumber', headerName: 'Start number', width: 200 },
                { field: 'fullName', headerName: 'Full name', width: 200 },
                { field: 'time', headerName: 'Corridor time', width: 200 },
                {
                    field: 'action', headerName: '', width: 200, disableClickEventBubbling: true,
                    renderCell: (params) => {
                        const onClick = () => {
                            const api = params.api;
                            const fields = api
                                .getAllColumns()
                                .map((c) => c.field)
                                .filter((c) => c !== "__check__" && !!c);
                            const thisRow = {};

                            fields.forEach((f) => {
                                thisRow[f] = params.getValue(f);
                            });
                            this.finishRunnerTime(thisRow['id']);
                        };

                        return <Button onClick={onClick}>Finish run</Button>;
                    }
                }
            ],
            runners: []
        }

        this.postRunnerTime = this.postRunnerTime.bind(this);
        this.getRunnerTimes = this.getRunnerTimes.bind(this);
        this.clearInputFields = this.clearInputFields.bind(this);
    }

    finishRunnerTime(id) {
        const runnerTime = {
            ID: id,
            FinishTime: new Date()
        }

        fetch('http://127.0.0.1:10000/updateRunnerTime', {
            method: 'UPDATE',
            headers: {
                'Content-Type': 'application/json',
            },
            body: JSON.stringify(runnerTime)
        }).then((result) => {
            this.getRunnerTimes();
        })
    }

    componentDidMount() {
        this.getRunnerTimes();
    }

    getRunnerTimes() {
        fetch('http://127.0.0.1:10000/getUnFinishedTimes')
            .then(res => res.json())
            .then((result) => {
                if (result != null) {
                    this.setState({ times: result });
                }
            })
    }

    onChange = event => {
        this.setState({
            [event.target.name]: event.target.value
        });
    };

    clearInputFields() {
        this.setState({
            chipNumber: "",
            fullName: ""
        });
    }

    postRunnerTime() {
        const runnerTime = {
            CorridorTime: new Date()
        }

        fetch('http://127.0.0.1:10000/postRunnerTime', {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json',
            },
            body: JSON.stringify(runnerTime)
        }).then(() => {
            this.getRunnerTimes();
        })
    }

    render() {
        return (
            <div>
                <button onClick={this.postRunnerTime}>Add time</button>

                <div>
                    <div style={{ height: 500, width: '100%' }}>
                        <DataGrid rows={this.state.times} columns={this.state.columns} pageSize={10} />
                    </div>
                </div>
            </div>
        );
    }
}

export default SubmitRunnerTime;