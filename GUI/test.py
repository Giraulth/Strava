import plotly.graph_objs as go
import pandas as pd
import requests
import json

data = requests.get('http://localhost:8080/api/getKudosRanking').text
df = pd.DataFrame(json.loads(data))


table = go.Figure(
    data=[
        go.Table(
            header=dict(
                values=[
                    'Utilisateur',
                    'Ratio',
                    'Kudos received',
                    'Activities seen'],
                fill=dict(
                    color='#C2D4FF'),
                align=[
                    'left',
                    'center']),
            cells=dict(
                values=[
                    df['username'],
                    df['ratio'],
                    df['kudos_count'],
                    df['activity_seen_count']],
                fill=dict(
                    color='#F5F8FF'),
                align=[
                    'left',
                    'center']))])

table.update_layout(
    title='Ranking Table',
    font=dict(size=12),
    height=300,
    margin=dict(l=20, r=20, t=40, b=20),
)

table.show()
