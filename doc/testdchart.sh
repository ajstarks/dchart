#!/bin/sh
dchart code/AAPL.d | pdfdeck -stdout - > dchart01.pdf
dchart -csv -csvcol=Date,Close code/AAPL.csv | pdfdeck -stdout - > dchart02.pdf
dchart -frame=t -framecolor=blue code/AAPL.d | pdfdeck -stdout - > dchart03.pdf
dchart -bgcolor=black -lcolor=white -vcolor=orange code/AAPL.d | pdfdeck -stdout - > dchart04.pdf
dchart -datacond=150,200,orange code/AAPL.d | pdfdeck -stdout - > dchart05.pdf
dchart -xlabel=2 code/AAPL.d | pdfdeck -stdout - > dchart06.pdf
dchart -xstagger code/AAPL.d | pdfdeck -stdout - > dchart07.pdf
dchart -xlabrot=300 code/AAPL.d | pdfdeck -stdout - > dchart08.pdf
dchart -vcolor=white -valpos=b code/AAPL.d | pdfdeck -stdout - > dchart09.pdf
dchart -vcolor=white -valpos=m code/AAPL.d | pdfdeck -stdout - > dchart10.pdf
dchart -vcolor=gray code/AAPL.d | pdfdeck -stdout - > dchart11.pdf
dchart -xlabel=2 -left 30 -right 70 -top 70 -bottom 40 code/AAPL.d | pdfdeck -stdout - > dchart12.pdf
dchart -color gray code/AAPL.d | pdfdeck -stdout - > dchart13.pdf
dchart -hline=250,Target -yaxis code/AAPL.d | pdfdeck -stdout - > dchart14.pdf
dchart -grid -yaxis code/AAPL.d | pdfdeck -stdout - > dchart15.pdf
dchart -yrange=0,300,25 -grid -yaxis code/AAPL.d | pdfdeck -stdout - > dchart16.pdf
dchart -barwidth=1 code/AAPL.d | pdfdeck -stdout - > dchart17.pdf
dchart -bar=f -dot code/AAPL.d | pdfdeck -stdout - > dchart18.pdf
dchart -bar=f -vol code/AAPL.d | pdfdeck -stdout - > dchart19.pdf
dchart -bar=f -vol -volop=90 code/AAPL.d | pdfdeck -stdout - > dchart20.pdf
dchart -bar=f -line code/AAPL.d | pdfdeck -stdout - > dchart21.pdf
dchart -bar=f -line -linewidth=0.5 code/AAPL.d | pdfdeck -stdout - > dchart22.pdf
dchart -bar=f -scatter code/AAPL.d | pdfdeck -stdout - > dchart23.pdf
dchart -bar=f -scatter -val=f code/AAPL.d | pdfdeck -stdout - > dchart24.pdf
dchart -bar=f -line -val=f -rline code/AAPL.d | pdfdeck -stdout - > dchart25.pdf
dchart -bar=f -line -val=f -rline -rlcolor=orange code/AAPL.d | pdfdeck -stdout - > dchart26.pdf
dchart -bar=f -line -vol -dot code/AAPL.d | pdfdeck -stdout - > dchart27.pdf
dchart -datafmt %0.3f -bar=f -dot -line code/AAPL.d | pdfdeck -stdout - > dchart28.pdf
dchart -bar=f -line -vol -dot -grid -yaxis code/AAPL.d | pdfdeck -stdout - > dchart29.pdf
dchart -hbar code/AAPL.d | pdfdeck -stdout - > dchart30.pdf
dchart -hbar -pct code/AAPL.d | pdfdeck -stdout - > dchart31.pdf
dchart -hbar -ls 1.5 code/AAPL.d | pdfdeck -stdout - > dchart32.pdf
dchart -wbar code/AAPL.d | pdfdeck -stdout - > dchart33.pdf
dchart -left=10 -right=25 -top=80 -bottom=60 -slope code/slope.d | pdfdeck -stdout - > dchart34.pdf
dchart -donut -color=std -pwidth=5 code/browser.d | pdfdeck -stdout - > dchart35.pdf
dchart -donut -color=std -title=f -top=70 -pwidth=20 -psize=20 code/browser.d | pdfdeck -stdout - > dchart36.pdf
dchart -pmap -pwidth=5 -textsize=1 code/browser.d | pdfdeck -stdout - > dchart37.pdf
dchart -pmap -pwidth=5 -textsize=1 -solidpmap code/browser.d | pdfdeck -stdout - > dchart38.pdf
dchart -pmap -pwidth=5 -textsize=1 -solidpmap -pmlen=30 code/browser.d | pdfdeck -stdout - > dchart39.pdf
dchart -left 35 -top 80 -ls 3 -pgrid -val=f code/incar.d | pdfdeck -stdout - > dchart40.pdf
dchart -radial -psize=10 -pwidth=20 -top=55 -textsize=3 code/count.d | pdfdeck -stdout - > dchart41.pdf
dchart -radial -psize=5 -pwidth=20 -top=55 -textsize=3 -spokes code/clock.d | pdfdeck -stdout - > dchart42.pdf
dchart -left 30 -top 80 -textsize 5 -lego  code/incar.d | pdfdeck -stdout -pagesize 1000,1000 - > dchart43.pdf
dchart -psize 30 -top=60 -fan code/occupation.d |  pdfdeck -stdout - > dchart44.pdf
dchart -psize 30 -top=60 -bowtie code/occupation.d |  pdfdeck -stdout - > dchart45.pdf
