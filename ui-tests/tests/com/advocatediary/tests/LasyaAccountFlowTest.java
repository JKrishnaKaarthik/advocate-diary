package com.advocatediary.tests;

import com.advocatediary.pages.DiaryPage;
import com.advocatediary.pages.SignInPage;
import org.testng.Assert;
import org.testng.annotations.BeforeMethod;
import org.testng.annotations.Test;

public class LasyaAccountFlowTest extends BaseTest {

  @Override
  @BeforeMethod
  public void openBrowser() {
    System.setProperty("baseUrl", "http://localhost:8080");
    System.setProperty("headless", "false");
    super.openBrowser();
  }

  @Test
  public void openSignInThenSignInShowsCourtDiary() {
    SignInPage signIn = new SignInPage(driver).open(baseUrl);
    Assert.assertTrue(signIn.currentUrl().contains("/sign-in"), "should start on /sign-in");
    Assert.assertEquals(signIn.heading(), "Welcome back");

    DiaryPage diary = signIn.signIn(LAWYER_EMAIL, LAWYER_PASSWORD);
    Assert.assertEquals(diary.heading(), "Court diary");
  }
}
