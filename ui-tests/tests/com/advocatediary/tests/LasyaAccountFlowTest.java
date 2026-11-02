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
    Assert.assertTrue(signIn.isOnSignIn(), "should start on /sign-in");

    DiaryPage diary = signIn.signIn(LAWYER_EMAIL, LAWYER_PASSWORD);
    Assert.assertTrue(diary.isOpen());
  }
}
