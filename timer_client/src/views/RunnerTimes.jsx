import React, { Component } from 'react';
import { DataGrid } from '@material-ui/data-grid';

class RunnerTimes extends Component {
    constructor(props) {
        super(props);
        this.state = {
            times: [],
            columns: [
                { field: 'id', headerName: 'Chip number', width: 200 },
                { field: 'startNumber', headerName: 'Start number', width: 200 },
                { field: 'fullName', headerName: 'Full name', width: 200 },
                { field: 'time', headerName: 'Finish time', width: 200 },
                { field: 'timingPoint', headerName: 'Timing point', width: 200 }
            ],
        }

        this.getRunnerTimes = this.getRunnerTimes.bind(this);
        this.addTimeToList = this.addTimeToList.bind(this);
    }

    componentDidMount() {
        this.getRunnerTimes();
    }

    getRunnerTimes() {
        fetch('http://127.0.0.1:10000/getAllRunnerTimes')
            .then(res => res.json())
            .then((result) => {
                if (result != null) {
                    this.setState({ times: result });
                }
            })
    }

    addTimeToList(time) {
        let timesHelper = [];
        this.setState({ times: timesHelper });
    }

    render() {
        return (
            <div>
                <div style={{ height: 500, width: '100%' }}>
                    <DataGrid rows={this.state.times} columns={this.state.columns} pageSize={10}
                        sortModel={[
                            {
                                field: 'time',
                                sort: 'desc',
                            },
                        ]}
                    />
                </div>
            </div>
        );
    }
}

export default RunnerTimes;