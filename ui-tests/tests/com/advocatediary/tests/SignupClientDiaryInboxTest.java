package com.advocatediary.tests;

import com.advocatediary.pages.CasePage;
import com.advocatediary.pages.ClientPage;
import com.advocatediary.pages.DiaryPage;
import com.advocatediary.pages.InboxPage;
import org.testng.Assert;
import org.testng.annotations.BeforeMethod;
import org.testng.annotations.Test;

public class SignupClientDiaryInboxTest extends BaseTest {

  @Override
  @BeforeMethod
  public void openBrowser() {
    System.setProperty("baseUrl", "http://localhost:8080");
    System.setProperty("headless", "false");
    super.openBrowser();
  }

  @Test
  public void signupAddClientCaseThenDiaryAndInboxShowHearing() {
    DiaryPage diary = signUpNewAccount();
    diary.waitUntilEmpty();
    Assert.assertTrue(diary.isEmpty(), "new account diary starts empty");

    ClientPage client = addClient(CLIENT_NAME);
    Assert.assertEquals(client.heading(), CLIENT_NAME);

    CasePage matter = openAndSaveCase(client);
    saveHearingOnCase(matter);

    diary = matter.nav().openDiary();
    diary.waitForSavedCase(CLIENT_NAME, CASE_TITLE, CASE_NUMBER, COURT, PURPOSE);
    Assert.assertTrue(diary.showsSavedCase(CLIENT_NAME, CASE_TITLE, CASE_NUMBER, COURT, PURPOSE));

    InboxPage inbox = diary.nav().openInbox();
    inbox.waitForNextHearing(CLIENT_NAME, CASE_TITLE, CASE_NUMBER);
    Assert.assertTrue(inbox.showsNextHearing(CLIENT_NAME, CASE_TITLE, CASE_NUMBER));
  }
}
