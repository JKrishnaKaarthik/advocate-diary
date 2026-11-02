package com.advocatediary.tests;

import com.advocatediary.pages.ClientsPage;
import com.advocatediary.pages.DiaryPage;
import org.testng.Assert;
import org.testng.annotations.Test;

public class AccountIsolationTest extends BaseTest {

  @Test
  public void secondAccountKeepsClientsAndDiarySeparate() {
    signupEmptyDiary();
    Assert.assertTrue(signOut(completeClientCaseOnDiary()).isOnSignIn());

    DiaryPage secondDiary = signupEmptyDiary();
    Assert.assertTrue(secondDiary.isEmpty(), "second account must not see the first account diary");

    ClientsPage clients = secondDiary.nav().openClients();
    clients.waitUntilEmpty();
    Assert.assertTrue(clients.isEmpty(), "second account must start with no clients");
    Assert.assertFalse(clients.hasClient(CLIENT_NAME), "second account must not see the first account client");

    DiaryPage listed = completeClientCaseOnDiary();
    Assert.assertTrue(diaryShowsMatter(listed));
    Assert.assertTrue(signOut(listed).isOnSignIn());
  }
}
