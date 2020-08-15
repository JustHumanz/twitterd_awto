from selenium import webdriver
from selenium.webdriver.common.keys import Keys
from selenium.webdriver.firefox.options import Options
from db import mydb,dbconn
import time,os,logging,random,argparse

logging.basicConfig(format='%(asctime)s %(message)s', datefmt='%m/%d/%Y %I:%M:%S %p')
logger=logging.getLogger()
logger.setLevel(logging.INFO)
#remove this if dev mode
options = Options()
options.add_argument('--headless')
options.add_argument('--hide-scrollbars')
options.add_argument('--disable-gpu')
class TwitterByURL:
    def __init__(self, username, password):
        self.username = username
        self.password = password
        self.bot = webdriver.Firefox(firefox_options = options)

    def login(self):
        bot = self.bot
        bot.get('https://twitter.com/login')
        logging.info("Url "+bot.current_url) 
        time.sleep(3)  # wait for the page to load it's contents
        email = bot.find_element_by_name('session[username_or_email]')
        password = bot.find_element_by_name('session[password]')
        email.clear()
        password.clear()
        email.send_keys(self.username)
        logging.info("Send Username")
        password.send_keys(self.password)
        logging.info("Send Password")
        password.send_keys(Keys.RETURN)
        time.sleep(3)

    def likeTweet(self,twitter_list):
        bot = self.bot
        count = 0
        for j in twitter_list:
            logging.info("Going to "+j)
            if count == 19:
                ran = random.randint(1, 3)
                logging.info("Sleep for "+str(ran)+" min and reset count")
                time.sleep(ran * 60)
                count = 0
            bot.get(j)
            time.sleep(3)
            bot.execute_script("window.scrollBy(0,400)")
            try:
                final = bot.find_elements_by_xpath("/html/body/div/div/div/div[2]/main/div/div/div/div[1]/div/div[2]/div/section/div/div/div[1]/div/div/article/div/div/div/div[3]/div[5]/div[3]/div[@data-testid='like']")
                for i in range(len(final)):
                    final[i].click()
                logging.info("done like "+j+" count "+str(count))
            except:
                logging.info("some error,LMAO ")

            count+=1
            time.sleep(1)

    def XD(self):
        logging.info("Taks complete")
        self.bot.quit()


class Database:
    def __init__(self, dbconn):
        self.db = dbconn
    
    def gettweet(self,member):
        db = self.db
        db.execute("SELECT PermanentURL FROM Vtuber.Twitter inner join VtuberMember on VtuberMember.id = VtuberData.VtuberMember_id where VtuberMember.VtuberName_EN= %s ",(member,))
        res = dbconn.fetchall()
        tweetID = []
        for i in res:
            tweetID.append(i[0])

        return tweetID

def main():
    parser = argparse.ArgumentParser(prog="main",description='Vtuber auto like[twitter]')
    parser.add_argument('--username', metavar="Twitter Username",dest="TWuser",help="Twitter username/email/no phone",required=True)
    parser.add_argument('--password', metavar="Twitter Password",dest="TWpass",help="Twitter password",required=True)
    parser.add_argument('--name', metavar="Name of Vtuber[En]",dest="MemberName",help="Vtuber member name",required=True)
    args = parser.parse_args()
    
    member = args.MemberName
    logging.info("Vtuber Member "+member)
    twuser = args.TWuser
    logging.info("Twitter Username "+twuser)
    twpass = args.TWpass
    logging.info("Twitter Password "+twpass)

    #get PermanentURL (example permanenURL : https://twitter.com/Aldin_Py/status/1286115777242750976 ) from DB
    tw_list = Database(dbconn).gettweet(member)
    logging.info("Get Twitter Permanen URL "+ str(len(tw_list)))

    #send User&pass
    brrr = TwitterByURL(twuser,twpass)
    #login 
    brrr.login() #takbiran ya balapan
    brrr.likeTweet(tw_list) #send like to permanentURL list
    brrr.XD()

main()