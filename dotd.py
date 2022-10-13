import re
import sys

import requests
import json
from libbgg.apiv1 import BGG
# You can also use version 2 of the api:
from libbgg.apiv2 import BGG as BGG2
from bs4 import BeautifulSoup
from libbgg.infodict import InfoDict
# import telegram


def fetchBGGLink(name):
    # name = "star wars: rebellion"
    stripped = re.sub(r'[^a-zA-Z0-9 ]', '', name).lower()
    stripped = stripped.replace("expansion", "")
    # print(stripped)
    conn = BGG2()

    results = conn.search(stripped)
    if int(results["items"]["total"]) > 0:
        x = results["items"]["item"]
        gameid = ""
        if isinstance(x, list):
            gameid = x[0].id
        if isinstance(x, InfoDict):
            gameid = x.get("id")
        return "https://boardgamegeek.com/boardgame/" + gameid
    return "no bgg link"
    # print(json.dumps(results, indent=4, sort_keys=False))
    # print(results["items"]["item"].get('id'))
    # for game in results["items"]["item"]:
    #     if (isinstance(game, )
    #     print(game)


# html = requests.get("https://www.gamenerdz.com/deal-of-the-day").text
# soup = BeautifulSoup(html, features="html.parser")
# card = [i for i in soup.find_all('article', 'card') if "Crisis" not in str(i)][0]
# title = card.find("h4", "card-title").a.text.strip()
# link = card.find("h4", "card-title").a.attrs['href']
# price = card.find("span", "price--withoutTax").text
#
# showPrice = True
# if "see price" in title:
#     showPrice = False
#
# title = title \
#     .replace("(Add to cart to see price)", "") \
#     .replace("(Deal of the Day)", "")
#
# old_title = title
#
# if showPrice:
#     title = f"[GN][DotD] - {title} - {price}"
# else:
#     title = f"[GN][DotD] - {title}"
#
# stock = BeautifulSoup(requests.get(link).text, features="html.parser").find("span", {"data-product-stock": True}).text

if len(sys.argv) != 2:
    sys.exit("must provide one arg")
bgglink = fetchBGGLink(sys.argv[1])
# out = title + "\n" + "in stock currently: " + stock + "\n" + bgglink
print(bgglink)


# bot = telegram.Bot(token="5715458728:AAHIcA3JGJiIWCq4l2NFPfbTKSDj586HmsQ")
#
# bot.sendMessage(504504957, out + "\n\n https://www.gamenerdz.com/deal-of-the-day" )


